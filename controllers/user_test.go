package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"ticatag_backend/db"
	"ticatag_backend/models"
	"ticatag_backend/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
