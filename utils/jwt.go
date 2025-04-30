package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


func GenerateJWT(userID string, role string, email string,createdAt int64) (string, error) {

	// Récupération de la clé secrète depuis les variables d'environnement
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("la variable d'environnement JWT_SECRET n'est pas définie")
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"email":   email,
		"createdAt":createdAt,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Expiration dans 24h
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string) (string, error) {

	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		fmt.Println("La variable JWT_SECRET n'est pas définie.")
		return "", errors.New("JWT secret not configured")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		//fmt.Println("La variable JWT_SECRET est publiée.")
		return []byte(secret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["user_id"].(string)
		//createdAt:= claims[]
		return userID, nil
	}

	return "", err
}
