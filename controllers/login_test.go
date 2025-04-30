package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"ticatag_backend/db"
	"ticatag_backend/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {
	// Setup Mongo de test
	SetupMongoTest()
	MongoDbTestConnection()

	// Créer un utilisateur avec mot de passe hashé
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := models.User{
		Username:  "testloginuser",
		Email:     "testlogin@example.com",
		Password:  string(hashedPassword),
		Role:      "user",
		CreatedAt: time.Now().Unix(),
	}
	collection := db.GetCollection("users")
	_, err := collection.InsertOne(context.TODO(), user)
	require.NoError(t, err)

	// Corps de la requête JSON
	loginInput := map[string]string{
		"username": "testloginuser",
		"password": "password123",
	}
	body, _ := json.Marshal(loginInput)

	// Configurer Gin
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/login", Login)

	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	token := resp["token"]
	assert.NotEmpty(t, token)

	//t.Logf("Token reçu : %s", token)
}

func TestRegister(t *testing.T) {
	// Préparer MongoDB de test
	SetupMongoTest()
	MongoDbTestConnection()

	// Nettoyer la collection users
	collection := db.GetCollection("users")
	_ = collection.Drop(context.TODO())

	// Corps JSON de la requête
	input := map[string]string{
		"username": "testuser_register",
		"email":    "register@example.com",
		"password": "pass123",
	}
	body, _ := json.Marshal(input)

	// Configurer Gin
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/register", Register)

	req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Vérifier la réponse
	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "User registered successfully", resp["message"])

	// Vérifier que l'utilisateur est bien en base
	var user models.User
	err = collection.FindOne(context.TODO(), bson.M{"username": "testuser_register"}).Decode(&user)
	require.NoError(t, err)
	assert.Equal(t, "register@example.com", user.Email)
	assert.NotEmpty(t, user.Password) // Doit être hashé
}
