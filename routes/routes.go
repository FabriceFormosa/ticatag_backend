package routes

import (
	"ticatag_backend/controllers"
	"ticatag_backend/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {

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

	protected := r.Group("/api/devices")

	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/profile", controllers.Profile)

		protected.GET("", controllers.GetDevices)
		protected.POST("", controllers.CreateDevice)
		protected.GET("/:id", controllers.GetDevice)
		protected.PUT("/:id", controllers.UpdateDevice)
		protected.DELETE("/:id", controllers.DeleteDevice)
		protected.GET("/search", controllers.FindDeviceByAddress)

	}

	return r

}
