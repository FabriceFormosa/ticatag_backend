package controllers

import (
	"context"
	"ticatag_backend/db"
	"ticatag_backend/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetDevices(c *gin.Context) {

	collection := db.DB.Collection("devices")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(500, gin.H{"error": "Erreur MongoDB"})
		return
	}
	defer cursor.Close(ctx)

	var devices []models.Device
	if err = cursor.All(ctx, &devices); err != nil {
		c.JSON(500, gin.H{"error": "Erreur decoding"})
		return
	}

	c.JSON(200, gin.H{"devices": devices})
}

func CreateDevice(c *gin.Context) {

	collection := db.DB.Collection("devices")

	var device models.Device
	if err := c.BindJSON(&device); err != nil {
		c.JSON(400, gin.H{"error": "Format invalide"})
		return
	}
	device.ID = primitive.NewObjectID()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, device)
	if err != nil {
		c.JSON(500, gin.H{"error": "Erreur lors de l'insertion"})
		return
	}
	c.JSON(201, device)
}

func GetDevice(c *gin.Context) {

	collection := db.DB.Collection("devices")

	idParam := c.Param("id")
	objID, _ := primitive.ObjectIDFromHex(idParam)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var device models.Device
	err := collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&device)
	if err != nil {
		c.JSON(404, gin.H{"error": "Utilisateur introuvable"})
		return
	}
	c.JSON(200, device)
}

func UpdateDevice(c *gin.Context) {

	collection := db.DB.Collection("devices")

	idParam := c.Param("id")
	objID, _ := primitive.ObjectIDFromHex(idParam)

	var device models.Device
	if err := c.BindJSON(&device); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{
		"$set": device,
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Utilisateur mis à jour"})
}

func DeleteDevice(c *gin.Context) {

	collection := db.DB.Collection("devices")

	idParam := c.Param("id")
	objID, _ := primitive.ObjectIDFromHex(idParam)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Utilisateur supprimé"})
}
