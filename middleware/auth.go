package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt" // vérifie ta version de jwt-go ou autre
)

//var secretKey = []byte("ton_secret")

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		fmt.Println("Appel fct AuthMiddleware ")

		// Récupération de la clé secrète depuis les variables d'environnement
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "la variable d'environnement JWT_SECRET n'est pas définie"})
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			fmt.Println("fct AuthMiddleware Authorization header missing")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			fmt.Println("AuthMiddleware :La variable secret est publiée ")
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			fmt.Println("Invalid token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Si tu veux récupérer des claims (infos dans le token)
		//claims := token.Claims.(jwt.MapClaims)

		c.Next()
	}
}
