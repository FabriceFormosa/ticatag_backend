package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"ticatag_backend/db"
	"ticatag_backend/middleware"
	"ticatag_backend/models"
	"ticatag_backend/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestProfile(t *testing.T) {
	// Setup MongoDB de test
	SetupMongoTest()
	MongoDbTestConnection()

	// Crée un utilisateur à insérer dans la base
	user := models.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		Role:      "user",
		CreatedAt: time.Now().Unix(),
	}
	collection := db.DB.Collection("users")
	res, err := collection.InsertOne(context.TODO(), user)
	require.NoError(t, err)
	insertedID := res.InsertedID.(primitive.ObjectID)

	// Génère un token JWT avec l'ID de l'utilisateur
	os.Setenv("JWT_SECRET", "test-secret")
	token, err := utils.GenerateJWT(insertedID.Hex(), user.Role, user.Email, user.CreatedAt) // à adapter selon ta fonction utils
	require.NoError(t, err)

	// Prépare le routeur
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(middleware.AuthMiddleware())
	router.GET("/profile", Profile)

	// Crée une requête avec le token dans l'en-tête Authorization
	req := httptest.NewRequest("GET", "/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	// Exécute la requête
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Vérifie la réponse
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, user.Username, response["username"])
	assert.Equal(t, user.Email, response["email"])
	assert.Equal(t, user.Role, response["role"])
	assert.Equal(t, insertedID.Hex(), response["id"])

	//t.Log("TestProfile terminé avec succès")
	//fmt.Println("TestProfile terminé avec succès")
}

func TestProfile_InvalidUserId(t *testing.T) {

	SetupMongoTest()
	MongoDbTestConnection()

	os.Setenv("JWT_SECRET", "test-secret")
	token, err := GenerateTestJWT()
	if err != nil {

		t.Log("Could not generate token")
	}

	// Crée un utilisateur à insérer dans la base
	user := models.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		Role:      "user",
		CreatedAt: time.Now().Unix(),
	}
	collection := db.DB.Collection("users")
	res, err := collection.InsertOne(context.TODO(), user)
	require.NoError(t, err)
	insertedID := res.InsertedID.(primitive.ObjectID)

	// Génère un token JWT avec l'ID de l'utilisateur
	os.Setenv("JWT_SECRET", "test-secret")
	_, err = utils.GenerateJWT(insertedID.Hex(), user.Role, user.Email, user.CreatedAt) // à adapter selon ta fonction utils
	require.NoError(t, err)

	// Corps de la requête JSON
	loginInput := map[string]string{
		"username": "testloginuser",
		"password": "password123",
	}
	body, _ := json.Marshal(loginInput)

	// Prépare le routeur
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(middleware.AuthMiddleware())
	router.GET("/profile", Profile)

	// Crée une requête avec le token dans l'en-tête Authorization
	req := httptest.NewRequest("GET", "/profile", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)

	// Exécute la requête
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid user ID")
}

func TestProfile_UserNotFound(t *testing.T) {
	SetupMongoTest()
	MongoDbTestConnection()

	// Génère un faux ID (non présent en base)
	fakeID := primitive.NewObjectID()

	// Génère un token JWT avec cet ID inexistant
	os.Setenv("JWT_SECRET", "test-secret")
	token, err := utils.GenerateJWT(fakeID.Hex(), "user", "ghost@example.com", time.Now().Unix())
	require.NoError(t, err)

	// Prépare le routeur
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(middleware.AuthMiddleware())
	router.GET("/profile", Profile)

	// Crée une requête GET avec le token dans l'en-tête Authorization
	req := httptest.NewRequest("GET", "/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	// Exécute la requête
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "User not found")
}

func TestProfile_MissingAuthorizationHeader(t *testing.T) {
	// Setup MongoDB de test
	SetupMongoTest()
	MongoDbTestConnection()

	// Crée un utilisateur à insérer dans la base
	user := models.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		Role:      "user",
		CreatedAt: time.Now().Unix(),
	}
	collection := db.DB.Collection("users")
	res, err := collection.InsertOne(context.TODO(), user)
	require.NoError(t, err)
	insertedID := res.InsertedID.(primitive.ObjectID)

	// Génère un token JWT avec l'ID de l'utilisateur
	os.Setenv("JWT_SECRET", "test-secret")
	_, err = utils.GenerateJWT(insertedID.Hex(), user.Role, user.Email, user.CreatedAt) // à adapter selon ta fonction utils
	require.NoError(t, err)

	// Prépare le routeur
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(middleware.AuthMiddleware())
	router.GET("/profile", Profile)

	// Crée une requête avec le token dans l'en-tête Authorization
	req := httptest.NewRequest("GET", "/profile", nil)

	// Exécute la requête
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authorization header missing")
}

func TestGetUsers(t *testing.T) {
	// Setup MongoDB de test
	SetupMongoTest()
	MongoDbTestConnection()

	// Initialiser JWT secret
	os.Setenv("JWT_SECRET", "test-secret")

	token, err := GenerateTestJWT()
	if err != nil {

		t.Log("Could not generate token")
	}

	// Insertion de quelques users pour le test
	collection := db.DB.Collection("users")
	usersToInsert := []any{
		models.User{Username: "User A", Email: "userA@gmail.com", Role: "user"},
		models.User{Username: "User B", Email: "userB@gmail.com0", Role: "user"},
	}
	res, err := collection.InsertMany(context.TODO(), usersToInsert)
	require.NoError(t, err)

	// Nettoyage après test
	t.Cleanup(func() {
		ids := make([]primitive.ObjectID, len(res.InsertedIDs))
		for i, id := range res.InsertedIDs {
			ids[i] = id.(primitive.ObjectID)
		}
		_, _ = collection.DeleteMany(context.TODO(), bson.M{"_id": bson.M{"$in": ids}})
	})

	// Setup router avec middleware JWT
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(middleware.AuthMiddleware())
	router.GET("/users", GetUsers)

	// Création de la requête GET
	req := httptest.NewRequest("GET", "/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Analyse de la réponse
	var response struct {
		Users []models.User `json:"users"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.GreaterOrEqual(t, len(response.Users), 2)
}
