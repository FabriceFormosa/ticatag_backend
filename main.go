package main

import (
	//"ticatag_backend/handlers"

	//"log"

	//"fmt"
	//"context"
	//"encoding/json"
	//"fmt"

	"ticatag_backend/db"
	"ticatag_backend/middleware"

	//"ticatag_backend/middleware"
	"time"

	//"ticatag_backend/models"

	//"ticatag_backend/routes"

	//"github.com/gin-gonic/gin"
	//"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/v2/bson"
	//"go.mongodb.org/mongo-driver/v2/mongo"

	"ticatag_backend/controllers"
	//"ticatag_backend/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	db.Connect()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Routes publiques
	r.POST("/login", controllers.Login)
	r.POST("/register", controllers.Register)
	r.GET("/api/devices/profile", controllers.Profile)

	protected := r.Group("/api/devices")

	protected.Use(middleware.AuthMiddleware())
	{

		protected.GET("", controllers.GetDevices)
		protected.POST("", controllers.CreateDevice)
		protected.GET("/:id", controllers.GetDevice)
		protected.PUT("/:id", controllers.UpdateDevice)
		protected.DELETE("/:id", controllers.DeleteDevice)
		protected.GET("/search", controllers.FindDeviceByAdress)

	}

	/* r.POST("/api/login", controllers.Login)
	protected := r.Group("/api")
	protected.Use(middleware.AuthRequired)
	protected.GET("/users", controllers.GetUsers) */

	//handlers.GetBooks(c *gi)

	//routes.SetupRoutesDevices(r)

	r.Run(":8080")
}
