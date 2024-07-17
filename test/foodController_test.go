package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	controllers "product-app/controller"
	"product-app/models"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
func TestCreateFood(t *testing.T){
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(),clientOptions)
	if err != nil{
		t.Fatalf("MongoDB connection error: %v", err)
	}
	defer client.Disconnect(context.TODO())
	
	testDatabase := client.Database("Library")
	userCollection := testDatabase.Collection("user")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = userCollection.DeleteMany(ctx, bson.M{})
	if err != nil {
		t.Fatalf("Error clearing test collection: %v", err)
	}

	gin.SetMode(gin.TestMode)
	router :=gin.Default()
	router.POST("/createFood",controllers.CreateFood)

	//Test verisi
	foodData := models.Food{
		Name :"muz",
		Price: 150,
	}
	jsonData ,_ := json.Marshal(foodData)

	req,_:= http.NewRequest(http.MethodPost,"/createFood",bytes.NewBuffer(jsonData))
	
	req.Header.Set("Content-Type", "application/json")
	// HTTP cevabını yakala
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Cevap durum kodunu kontrol et
	assert.Equal(t, http.StatusCreated, w.Code)
	var responseFood models.Food
	err = json.Unmarshal(w.Body.Bytes(), &responseFood)
	if err != nil {
		t.Fatalf("Error unmarshaling response: %v", err)
	}
	assert.Equal(t, responseFood.Name, responseFood.Name)
	assert.Equal(t, responseFood.Price, responseFood.Price)
	assert.NotEmpty(t, responseFood.ID)
}