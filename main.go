package main

import (
	//"ticatag_backend/handlers"

	//"log"

	//"fmt"
	//"context"
	//"encoding/json"
	//"fmt"
	"ticatag_backend/db"
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

	//db.ConnectDB() // Connexion MongoDB

	db.Connect()

	//var coll = config.GetCollection("movies")cls

	//GetBookByTitle("Millénium")

	//GetBooks()

	//CreateOneBook("Atonement", "Ian McEwan")

	r := gin.Default()

	//r.Use(cors.Default())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/api/devices", controllers.CreateDevice)
	r.GET("/api/devices", controllers.GetDevices)
	r.GET("/api/devices/:id", controllers.GetDevice)
	r.PUT("/api/devices/:id", controllers.UpdateDevice)
	r.DELETE("/api/devices/:id", controllers.DeleteDevice)

	/* r.POST("/api/login", controllers.Login)
	protected := r.Group("/api")
	protected.Use(middleware.AuthRequired)
	protected.GET("/users", controllers.GetUsers) */

	//handlers.GetBooks(c *gi)

	//routes.SetupRoutesDevices(r)

	r.Run(":8080")
}

/*
func GetBookByTitle(title string) {

	var coll = config.GetCollection("books")

	var result bson.M
	err := coll.FindOne(context.TODO(), bson.D{{Key: "title", Value: title}}).
		Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the title %s\n", title)
		return
	}
	if err != nil {
		panic(err)
	}
	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonData)

}

func GetBooks() {

	var coll = config.GetCollection("books")

	// Lecture de tous les éléments
	cursor, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}

	for cursor.Next(context.TODO()) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("%+v\n", result)
		jsonData, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", jsonData)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

}

func CreateOneBook(new_title string, new_author string) {

	var coll = config.GetCollection("books")

	doc := models.Book{Title: new_title, Author: new_author}

	result, err := coll.InsertOne(context.TODO(), doc)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
}
*/
