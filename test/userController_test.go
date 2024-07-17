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

// TestSignUp is a test function for the SignUp endpoint
func TestSignUp(t *testing.T) {
	// MongoDB bağlantısını oluştur
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		t.Fatalf("MongoDB connection error: %v", err)
	}
	defer client.Disconnect(context.TODO())

	// Veritabanını ve koleksiyonu seç
	testDatabase := client.Database("Library")

	userCollection := testDatabase.Collection("user")

	// Test veritabanını temizle
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = userCollection.DeleteMany(ctx, bson.M{})
	if err != nil {
		t.Fatalf("Error clearing test collection: %v", err)
	}

	// Test router'ını oluştur
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/signup", controllers.SignUp)

	// Test verisi
	userData := models.User{
		Email:    "test@gmail.com",
		Password: "password123",
	}
	jsonData, _ := json.Marshal(userData)

	// HTTP isteği simüle et
	req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// HTTP cevabını yakala
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Cevap durum kodunu kontrol et
	assert.Equal(t, http.StatusCreated, w.Code)

	// Cevap gövdesini kontrol et
	var responseUser models.User
	err = json.Unmarshal(w.Body.Bytes(), &responseUser)
	if err != nil {
		t.Fatalf("Error unmarshaling response: %v", err)
	}
	assert.Equal(t, userData.Email, responseUser.Email)

	// Veritabanında kullanıcının varlığını kontrol et
	err = userCollection.FindOne(ctx, bson.M{"email": userData.Email}).Decode(&responseUser)
	assert.Nil(t, err)
	assert.Equal(t, userData.Email, responseUser.Email)
	assert.NotEmpty(t, responseUser.ID)
	assert.NotEqual(t, userData.Password, responseUser.Password) // Şifre hashlenmiş olmalı
}
