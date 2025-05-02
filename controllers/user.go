package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	//"net/http"
	"ticatag_backend/db"
	"ticatag_backend/models"
	"ticatag_backend/utils"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUsers(c *gin.Context) {
	collection := db.DB.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(500, gin.H{"error": "Erreur MongoDB"})
		return
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		c.JSON(500, gin.H{"error": "Erreur decoding"})
		return
	}

	c.JSON(200, gin.H{"users": users})
}

/*
Lit le token JWT dans l'Authorization Header

Vérifie l'authentification

Récupère l'utilisateur en base MongoDB

Renvoie les données du profil
*/

// Fonction pour obtenir le profil utilisateur
func Profile(c *gin.Context) {
	fmt.Println("Appel fct Profile ")
	// Récupère le token Authorization: Bearer <token>
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		return
	}

	// Extrait le token
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		return
	}

	// Vérifie le token et récupère l'ID utilisateur
	userID, err := utils.ParseToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Convertit l'ID en ObjectID MongoDB
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Récupère l'utilisateur dans MongoDB
	userCollection := db.GetCollection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err = userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Réponse propre (on ne renvoie pas le mot de passe !)
	c.JSON(http.StatusOK, gin.H{
		"id":        user.ID.Hex(),
		"username":  user.Username,
		"email":     user.Email,
		"role":      user.Role,
		"createdAt": user.CreatedAt,
	})
}
