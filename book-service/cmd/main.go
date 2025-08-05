package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"book-service/internal/handler"
	"book-service/internal/repository"
	"book-service/internal/routes"
	"book-service/internal/service"

	_ "book-service/docs"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @title Booktopia Book Service API
// @version 1.0
// @description Ini adalah service Manajemen Book untuk Booktopia.
// @termsOfService http://swagger.io/terms/

// @contact.name Tim Developer Booktopia
// @contact.email dev@booktopia.com

// @host localhost:8081
// @BasePath /api
func main() {
	// 1. Muat variabel dari file .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Could not find or load .env file, using system environment variables")
	}

	// 2. Baca konfigurasi dari environment variable
	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB")
	port := os.Getenv("PORT") // PORT untuk server HTTP, bukan GRPC

	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable is not set")
	}
	if dbName == "" {
		log.Fatal("MONGO_DB environment variable is not set")
	}
	if port == "" {
		port = "8081" // Gunakan port default jika tidak diset
	}

	// 3. Konfigurasi Koneksi MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	bookCollection := client.Database(dbName).Collection("books")

	// 4. Inisialisasi Layer (Dependency Injection)
	bookRepo := repository.NewBookRepository(bookCollection)
	bookService := service.NewBookService(bookRepo)
	bookHandler := handler.NewBookHandler(bookService)

	// 5. Setup HTTP Server & Routing
	e := echo.New()
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// 6. Setup Route
	routes.SetupRoutes(e, bookHandler)

	// 7. Jalankan Server
	serverPort := ":" + port
	fmt.Printf("Book Service with MongoDB is running on port %s\n", serverPort)
	e.Logger.Fatal(e.Start(serverPort))
}