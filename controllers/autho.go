package controllers

import (
	"context"
	"net/http"
	"ticatag_backend/db"
	"ticatag_backend/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
)

var secretKey = []byte("supersecret")

func Login(c *gin.Context) {
	var input struct {
		Email string `json:"email"`
	}
	c.BindJSON(&input)

	user := models.User{Email: input.Email}

	collection := db.DB.Collection("users")

	// Vérifie si l'utilisateur existe
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"email": input.Email}).Err()
	if err != nil {
		// Si non trouvé → on l'ajoute
		_, _ = collection.InsertOne(ctx, user)
	}

	// Création du token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, _ := token.SignedString(secretKey)

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}
