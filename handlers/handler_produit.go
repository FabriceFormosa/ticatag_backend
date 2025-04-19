package handlers

// "context"
//"ticatag_backend/db"
/* 	"ticatag_backend/models"
"net/http"
"time"
"github.com/gin-gonic/gin"
"go.mongodb.org/mongo-driver/bson"
"go.mongodb.org/mongo-driver/bson/primitive"
*/

//var collection = db.GetCollection("produits")

/* func GetProduits(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var produits []models.Produit
	if err := cursor.All(ctx, &produits); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, produits)
} */
/*
func CreateProduit(c *gin.Context) {
	var produit models.Produit
	if err := c.ShouldBindJSON(&produit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	produit.ID = primitive.NewObjectID()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, produit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, produit)
}

func GetProduitByID(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var produit models.Produit
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&produit)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produit non trouvé"})
		return
	}

	c.JSON(http.StatusOK, produit)
}
func UpdateProduit(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	var data models.Produit
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"nom":  data.Nom,
			"prix": data.Prix,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil || result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produit non trouvé ou échec de la mise à jour"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Produit mis à jour"})
}
func DeleteProduit(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil || result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produit non trouvé"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Produit supprimé"})
}
*/

// (Ajoute GetProduitByID, UpdateProduit, DeleteProduit si tu veux la suite)
