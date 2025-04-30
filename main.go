package main

import (
	"ticatag_backend/db"
	"ticatag_backend/routes"
)

func main() {

	db.Connect()

	router := routes.SetupRoutes()

	router.Run(":8080")
}
