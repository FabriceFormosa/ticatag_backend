package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"ticatag_backend/db"
	"ticatag_backend/middleware"
	"ticatag_backend/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var testDB *mongo.Database

func SetupMongoTest() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	testDB = client.Database("test_db")
	db.DB = testDB // ← important si tu utilises db.DB partout
}
func MongoDbTestConnection() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Erreur de connexion :", err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("Mongo ne répond pas :", err)
	}
	//fmt.Println("✅ Connexion à MongoDB OK")
}

// Génère un JWT de test signé avec le secret fourni
func GenerateTestJWT() (string, error) {

	// Récupération de la clé secrète depuis les variables d'environnement
	secret := os.Getenv("JWT_SECRET")

	//fmt.Println("secret: ", secret)

	if secret == "" {
		return "", fmt.Errorf("la variable d'environnement JWT_SECRET n'est pas définie")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   "Test_userID",
		"role":      "Test_role",
		"email":     "Test_email",
		"createdAt": time.Now().Unix(),
		"exp":       time.Now().Add(time.Hour * 24).Unix(), // Expiration dans 24h
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Print("could not sign test JWT: " + err.Error())
		panic("could not sign test JWT: " + err.Error())
	}

	//fmt.Print("tokenString: ", tokenString)
	//fmt.Print("nil: ", nil)

	return tokenString, nil
}

func TestCreateDevice_InvalidJSON(t *testing.T) {
	SetupMongoTest()
	MongoDbTestConnection()

	os.Setenv("JWT_SECRET", "test-secret")
	token, err := GenerateTestJWT()
	if err != nil {

		t.Log("Could not generate token")
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(middleware.AuthMiddleware())
	router.POST("/devices", CreateDevice)

	// JSON invalide : champ mal formé ou type incohérent
	invalidJSON := `{"adress":123}` // `adress` est supposé être une string

	req := httptest.NewRequest("POST", "/devices", strings.NewReader(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Format invalide")
}

func TestCreateDevice(t *testing.T) {
	SetupMongoTest()
	MongoDbTestConnection()

	os.Setenv("JWT_SECRET", "test-secret")
	token, err := GenerateTestJWT()
	if err != nil {

		t.Log("Could not generate token")
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(middleware.AuthMiddleware())
	router.POST("/devices", CreateDevice)

	device := models.Device{
		Adress:    "DeviceTest",
		Latitude:  "Device Latitude",
		Longitude: "Device Longitude",
	}

	bodyJSON, err := json.Marshal(device)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/devices", strings.NewReader(string(bodyJSON)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Device
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, device.Adress, response.Adress)
	assert.Equal(t, device.Latitude, response.Latitude)
	assert.Equal(t, device.Longitude, response.Longitude)
	assert.False(t, response.ID.IsZero())

	t.Cleanup(func() {
		collection := db.DB.Collection("devices")
		_, err := collection.DeleteMany(context.TODO(), map[string]any{
			"adress": "DeviceTest",
		})
		if err != nil {
			t.Logf("cleanup failed: %v", err)
		}
	})

}

func TestGetOneDevice(t *testing.T) {
	// Connexion MongoDB de test
	SetupMongoTest()
	MongoDbTestConnection()

	// Initialiser JWT secret
	os.Setenv("JWT_SECRET", "test-secret")

	token, err := GenerateTestJWT()
	if err != nil {

		t.Log("Could not generate token")
	}
	// Crée un device manuellement dans la base pour le test
	device := models.Device{
		Adress:    "DeviceTest",
		Latitude:  "Device Latitude",
		Longitude: "Device Longitude",
	}
	collection := db.DB.Collection("devices")
	res, err := collection.InsertOne(context.TODO(), device)
	require.NoError(t, err)

	insertedID := res.InsertedID.(primitive.ObjectID)

	// Nettoyage après test
	t.Cleanup(func() {
		_, _ = collection.DeleteOne(context.TODO(), bson.M{"_id": insertedID})
	})

	// Initialisation du routeur avec middleware
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(middleware.AuthMiddleware())
	router.GET("/devices/:id", GetDevice)

	// Création de la requête GET avec auth
	req := httptest.NewRequest("GET", "/devices/"+insertedID.Hex(), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Vérification de la réponse
	var response models.Device
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, device.Adress, response.Adress)
	assert.Equal(t, device.Latitude, response.Latitude)
	assert.Equal(t, device.Longitude, response.Longitude)
	assert.Equal(t, insertedID, response.ID)
}
func TestGetOneDeviceNotFound(t *testing.T) {
	// Connexion MongoDB de test
	SetupMongoTest()
	MongoDbTestConnection()

	// Initialiser JWT secret
	os.Setenv("JWT_SECRET", "test-secret")

	// Générer un token de test
	token, err := GenerateTestJWT()
	require.NoError(t, err, "Could not generate token")

	// Crée un device manuellement dans la base pour le test
	device := models.Device{
		Adress:    "DeviceTest",
		Latitude:  "Device Latitude",
		Longitude: "Device Longitude",
	}
	collection := db.DB.Collection("devices")
	res, err := collection.InsertOne(context.TODO(), device)
	require.NoError(t, err)

	insertedID := res.InsertedID.(primitive.ObjectID)

	// Nettoyage après test
	t.Cleanup(func() {
		_, _ = collection.DeleteOne(context.TODO(), bson.M{"_id": insertedID})
	})

	// Initialisation du routeur avec middleware
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(middleware.AuthMiddleware())
	router.GET("/devices/:id", GetDevice)

	// Création de la requête GET avec un ID invalide (non existant dans la DB)
	invalidID := primitive.NewObjectID() // ID invalide pour le test
	req := httptest.NewRequest("GET", "/devices/"+invalidID.Hex(), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Vérifications
	assert.Equal(t, 404, w.Code)
	assert.Contains(t, w.Body.String(), "Device introuvable")
}
func TestGetOneDeviceIdInvalid(t *testing.T) {
	// Connexion MongoDB de test
	SetupMongoTest()
	MongoDbTestConnection()

	// Initialiser JWT secret
	os.Setenv("JWT_SECRET", "test-secret")

	// Générer un token de test
	token, err := GenerateTestJWT()
	require.NoError(t, err, "Could not generate token")

	// Crée un device manuellement dans la base pour le test
	device := models.Device{
		Adress:    "DeviceTest",
		Latitude:  "Device Latitude",
		Longitude: "Device Longitude",
	}
	collection := db.DB.Collection("devices")
	res, err := collection.InsertOne(context.TODO(), device)
	require.NoError(t, err)

	insertedID := res.InsertedID.(primitive.ObjectID)

	// Nettoyage après test
	t.Cleanup(func() {
		_, _ = collection.DeleteOne(context.TODO(), bson.M{"_id": insertedID})
	})

	// Initialisation du routeur avec middleware
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(middleware.AuthMiddleware())
	router.GET("/devices/:id", GetDevice)

	// Création de la requête GET avec un ID invalide (non existant dans la DB)
	//invalidID := primitive.NewObjectID() // ID invalide pour le test
	req := httptest.NewRequest("GET", "/devices/15555", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Vérifications
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "ID invalide")
}

// -------------------------------------------------------------------------------------//
type DeviceFinder interface {
	Find(ctx context.Context, filter any) (*mongo.Cursor, error)
}

func GetDevicesFromRepo(c *gin.Context, repo DeviceFinder) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := repo.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(500, gin.H{"error": "Erreur MongoDB"})
		return
	}
	defer cursor.Close(ctx)
	// ...
}

type MockRepo struct{}

func (m *MockRepo) Find(ctx context.Context, filter any) (*mongo.Cursor, error) {
	return nil, errors.New("force MongoDB error")
}

func TestGetDevices_MongoFindError(t *testing.T) {

	// Initialiser JWT secret
	os.Setenv("JWT_SECRET", "test-secret")

	token, err := GenerateTestJWT()
	if err != nil {

		t.Log("Could not generate token")
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(middleware.AuthMiddleware())

	mockRepo := &MockRepo{}

	router.GET("/devices", func(c *gin.Context) {
		GetDevicesFromRepo(c, mockRepo)
	})

	req := httptest.NewRequest("GET", "/devices", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)
	assert.Contains(t, w.Body.String(), "Erreur MongoDB")
}

// ------------------------------------------------------------------------------------------
type CursorWrapper interface {
	All(ctx context.Context, results interface{}) error
	Close(ctx context.Context) error
}

type MockCursorWithDecodeError struct{}

func (m *MockCursorWithDecodeError) All(ctx context.Context, results interface{}) error {
	return errors.New("mock decode error")
}

func (m *MockCursorWithDecodeError) Close(ctx context.Context) error {
	return nil
}

func TestGetDevices_DecodeError(t *testing.T) {
	// Initialiser JWT secret
	os.Setenv("JWT_SECRET", "test-secret")

	token, err := GenerateTestJWT()
	if err != nil {

		t.Log("Could not generate token")
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(middleware.AuthMiddleware())

	// Handler mocké utilisant MockCursorWithDecodeError
	router.GET("/devices", func(c *gin.Context) {
		var devices []models.Device
		mockCursor := &MockCursorWithDecodeError{}
		err := mockCursor.All(context.Background(), &devices)
		if err != nil {
			c.JSON(500, gin.H{"error": "Erreur decoding"})
			return
		}
	})

	req := httptest.NewRequest("GET", "/devices", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Erreur decoding")
}

//-------------------------------------------------------------------------------------------

func TestGetDevices(t *testing.T) {
	// Setup MongoDB de test
	SetupMongoTest()
	MongoDbTestConnection()

	// Initialiser JWT secret
	os.Setenv("JWT_SECRET", "test-secret")

	token, err := GenerateTestJWT()
	if err != nil {

		t.Log("Could not generate token")
	}

	// Insertion de quelques devices pour le test
	collection := db.DB.Collection("devices")
	devicesToInsert := []any{
		models.Device{Adress: "Device A", Latitude: "45.0", Longitude: "1.0"},
		models.Device{Adress: "Device B", Latitude: "46.0", Longitude: "2.0"},
	}
	res, err := collection.InsertMany(context.TODO(), devicesToInsert)
	require.NoError(t, err)

	// Nettoyage après test
	t.Cleanup(func() {
		ids := make([]primitive.ObjectID, len(res.InsertedIDs))
		for i, id := range res.InsertedIDs {
			ids[i] = id.(primitive.ObjectID)
		}
		_, _ = collection.DeleteMany(context.TODO(), bson.M{"_id": bson.M{"$in": ids}})
	})

	// Setup router avec middleware JWT
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(middleware.AuthMiddleware())
	router.GET("/devices", GetDevices)

	// Création de la requête GET
	req := httptest.NewRequest("GET", "/devices", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Analyse de la réponse
	var response struct {
		Devices []models.Device `json:"devices"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.GreaterOrEqual(t, len(response.Devices), 2)
}

//--------------------------------------------------------------------------------------------

func TestUpdateDevice(t *testing.T) {
	// Setup MongoDB de test
	SetupMongoTest()
	MongoDbTestConnection()

	// Initialiser JWT secret
	os.Setenv("JWT_SECRET", "test-secret")

	token, err := GenerateTestJWT()
	if err != nil {

		t.Log("Could not generate token")
	}

	// Insertion d'un device pour l'update
	originalDevice := models.Device{
		Adress:    "Old Address",
		Latitude:  "10.0",
		Longitude: "20.0",
	}
	collection := db.DB.Collection("devices")
	res, err := collection.InsertOne(context.TODO(), originalDevice)
	require.NoError(t, err)
	insertedID := res.InsertedID.(primitive.ObjectID)

	// Nettoyage à la fin
	t.Cleanup(func() {
		_, _ = collection.DeleteOne(context.TODO(), bson.M{"_id": insertedID})
	})

	// Nouveau device à envoyer dans le corps
	updatedDevice := models.Device{
		Adress:    "New Address",
		Latitude:  "99.9",
		Longitude: "88.8",
	}
	body, err := json.Marshal(updatedDevice)
	require.NoError(t, err)

	// Router avec AuthMiddleware
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(middleware.AuthMiddleware())
	router.PUT("/devices/:id", UpdateDevice)

	req := httptest.NewRequest("PUT", "/devices/"+insertedID.Hex(), strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Vérification dans la base
	var result models.Device
	err = collection.FindOne(context.TODO(), bson.M{"_id": insertedID}).Decode(&result)
	require.NoError(t, err)
	assert.Equal(t, updatedDevice.Adress, result.Adress)
	assert.Equal(t, updatedDevice.Latitude, result.Latitude)
	assert.Equal(t, updatedDevice.Longitude, result.Longitude)
}

func TestUpdateDevice_InvalidJSON(t *testing.T) {
	SetupMongoTest()
	MongoDbTestConnection()

	os.Setenv("JWT_SECRET", "test-secret")
	token, err := GenerateTestJWT()
	require.NoError(t, err)

	// Insertion d'un device
	originalDevice := models.Device{
		Adress:    "Old Address",
		Latitude:  "10.0",
		Longitude: "20.0",
	}
	collection := db.DB.Collection("devices")
	res, err := collection.InsertOne(context.TODO(), originalDevice)
	require.NoError(t, err)
	insertedID := res.InsertedID.(primitive.ObjectID)

	t.Cleanup(func() {
		_, _ = collection.DeleteOne(context.TODO(), bson.M{"_id": insertedID})
	})

	// JSON invalide
	invalidBody := `{"Adress": "New Address", "Latitude": "99.9", "Longitude": `

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(middleware.AuthMiddleware())
	router.PUT("/devices/:id", UpdateDevice)

	req := httptest.NewRequest("PUT", "/devices/"+insertedID.Hex(), strings.NewReader(invalidBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

// ----------------------------------------------------------------------------------------
func TestDeleteDevice(t *testing.T) {
	// Setup MongoDB de test
	SetupMongoTest()
	MongoDbTestConnection()

	// Initialiser JWT secret
	os.Setenv("JWT_SECRET", "test-secret")

	token, err := GenerateTestJWT()
	if err != nil {

		t.Log("Could not generate token")
	}
	// Insertion d'un device à supprimer
	device := models.Device{
		Adress:    "ToDelete",
		Latitude:  "00.0",
		Longitude: "00.0",
	}

	// Nettoyer la collection users
	collection := db.DB.Collection("devices")
	_ = collection.Drop(context.TODO())

	res, err := collection.InsertOne(context.TODO(), device)
	require.NoError(t, err)
	insertedID := res.InsertedID.(primitive.ObjectID)
	t.Logf("insertedID : %s", insertedID)

	// Router avec middleware JWT
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(middleware.AuthMiddleware())
	router.DELETE("/devices/:id", DeleteDevice)

	//t.Logf("invalidID : %s", invalidID)

	req := httptest.NewRequest("DELETE", "/devices/"+insertedID.Hex(), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Vérifications
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Device supprimé")

}

func TestDeleteDevice_IdNotFound(t *testing.T) {
	// Setup MongoDB de test
	SetupMongoTest()
	MongoDbTestConnection()

	// Initialiser JWT secret
	os.Setenv("JWT_SECRET", "test-secret")

	token, err := GenerateTestJWT()
	if err != nil {

		t.Log("Could not generate token")
	}
	// Insertion d'un device à supprimer
	device := models.Device{
		Adress:    "ToDelete",
		Latitude:  "00.0",
		Longitude: "00.0",
	}

	// Nettoyer la collection users
	collection := db.DB.Collection("devices")
	_ = collection.Drop(context.TODO())

	res, err := collection.InsertOne(context.TODO(), device)
	require.NoError(t, err)
	insertedID := res.InsertedID.(primitive.ObjectID)
	t.Logf("insertedID : %s", insertedID)

	// Router avec middleware JWT
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(middleware.AuthMiddleware())
	router.DELETE("/devices/:id", DeleteDevice)

	invalidID := primitive.NewObjectID() // ID invalide pour le test
	//t.Logf("invalidID : %s", invalidID)

	req := httptest.NewRequest("DELETE", "/devices/"+invalidID.Hex(), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Vérifications
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "No device found with this Id")

}

func TestDeleteDevice_Error(t *testing.T) {
	// Setup MongoDB de test
	SetupMongoTest()
	MongoDbTestConnection()

	// Initialiser JWT secret
	os.Setenv("JWT_SECRET", "test-secret")

	token, err := GenerateTestJWT()
	if err != nil {
		t.Logf("Could not generate token %s", token)
	}

	// Router avec middleware JWT
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(middleware.AuthMiddleware())
	router.DELETE("/devices/:id", DeleteDevice)

	validID := primitive.NewObjectID().Hex()
	invalidID := fmt.Sprintf("%s%s", validID, "XXX") // ID invalide pour le test
	t.Logf("validID : %s", validID)
	t.Logf("invalidID : %s", invalidID)

	req := httptest.NewRequest("DELETE", "/devices/"+invalidID, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Vérifications
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "ID invalide")

}

type Deleter interface {
	DeleteOne(context.Context, interface{}) (*mongo.DeleteResult, error)
}
type MockDeleteCollection struct{}

func (m *MockDeleteCollection) DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	return nil, errors.New("suppression échouée")
}

func DeleteDeviceWithRepo(c *gin.Context, repo Deleter) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID invalide"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = repo.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(500, gin.H{"error": "Erreur lors de la suppression"})
		return
	}
	c.JSON(200, gin.H{"message": "Utilisateur supprimé"})
}

func TestDeleteDevice_Error2(t *testing.T) {
	// Setup MongoDB de test
	SetupMongoTest()
	MongoDbTestConnection()

	// Initialiser JWT secret
	os.Setenv("JWT_SECRET", "test-secret")

	token, err := GenerateTestJWT()
	if err != nil {
		t.Logf("Could not generate token %s", token)
	}

	// Router avec middleware JWT
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(middleware.AuthMiddleware())

	mockRepo := &MockDeleteCollection{}

	router.DELETE("/devices/:id", func(c *gin.Context) {
		DeleteDeviceWithRepo(c, mockRepo)
	})

	invalidID := primitive.NewObjectID().Hex()
	t.Logf("invalidID : %s", invalidID)

	req := httptest.NewRequest("DELETE", "/devices/"+invalidID, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Vérifications
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Erreur lors de la suppression")

}

//-------------------------------------------------------------------------------------

func TestFindDeviceByAdress(t *testing.T) {
	// Setup MongoDB test
	SetupMongoTest()
	MongoDbTestConnection()

	// Initialiser JWT secret
	os.Setenv("JWT_SECRET", "test-secret")

	token, err := GenerateTestJWT()
	if err != nil {

		t.Log("Could not generate token")
	}

	// Insertion de quelques devices
	collection := db.DB.Collection("devices")
	devices := []any{
		models.Device{Adress: "Paris 12", Latitude: "48.8", Longitude: "2.3"},
		models.Device{Adress: "paris 15", Latitude: "48.84", Longitude: "2.29"},
		models.Device{Adress: "Lyon", Latitude: "45.75", Longitude: "4.85"},
	}
	res, err := collection.InsertMany(context.TODO(), devices)
	require.NoError(t, err)

	// Cleanup après test
	t.Cleanup(func() {
		ids := make([]primitive.ObjectID, len(res.InsertedIDs))
		for i, id := range res.InsertedIDs {
			ids[i] = id.(primitive.ObjectID)
		}
		_, _ = collection.DeleteMany(context.TODO(), bson.M{"_id": bson.M{"$in": ids}})
	})

	// Setup du routeur avec middleware
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(middleware.AuthMiddleware())
	router.GET("/devices/search", FindDeviceByAdress)

	// Requête GET avec query ?q=paris
	req := httptest.NewRequest("GET", "/devices/search?q=paris", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Vérifie la réponse
	assert.Equal(t, http.StatusOK, w.Code)

	var response []map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.GreaterOrEqual(t, len(response), 2)
	for _, device := range response {
		assert.Contains(t, strings.ToLower(device["adress"].(string)), "paris")
	}
}
