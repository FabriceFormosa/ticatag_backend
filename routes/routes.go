package routes

import (
	"ticatag_backend/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutesDevices(router *gin.Engine) {

	r := router.Group("/devices")
	r.POST("/", handlers.CreateDevice)

}

func SetupRoutes(router *gin.Engine) {

	/* r := router.Group("/books")
	{
		r.GET("/", handlers.GetBooks)
	} */

	/* 	r := router.Group("/produits")
		{

	 		r.GET("/", handlers.GetProduits)
	 		r.GET("/:id", handlers.GetProduitByID)
			r.POST("/", handlers.CreateProduit)
			r.PUT("/:id", handlers.UpdateProduit)
			r.DELETE("/:id", handlers.DeleteProduit)
		}  */
}
