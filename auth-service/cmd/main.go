package main

import (
	"auth-service/internal/handler"
	"auth-service/internal/model"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	"auth-service/pkg/hasher"
	"auth-service/pkg/jwt"
	"auth-service/pkg/mail"
	"fmt"
	"log"
	"os"

	_ "auth-service/docs"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Booktopia Auth Service API
// @version 1.0
// @description Ini adalah service autentikasi untuk Booktopia.
// @termsOfService http://swagger.io/terms/

// @contact.name Tim Developer Booktopia
// @contact.email dev@booktopia.com

// @host localhost:8082
// @BasePath /api/auth
func main() {
	// 1. Memuat .env langsung di sini
	godotenv.Load()
	// if err != nil {
	// 	log.Fatalf("Error loading .env file: %v", err)
	// }

	// 2. Inisialisasi Database langsung di sini
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	portDB := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, pass, name, portDB)
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Success connect to database")


	// Inisialisasi semua komponen Infrastruktur
	passwordHasher := hasher.NewBcrypt()
	jwtManager := jwt.NewManager(os.Getenv("JWT_SECRET"))
	mailer := mail.NewSmtpMailer(
		os.Getenv("MAILTRAP_HOST"),
		os.Getenv("MAILTRAP_PORT"),
		os.Getenv("MAILTRAP_USER"),
		os.Getenv("MAILTRAP_PASS"),
		os.Getenv("MAIL_SENDER"),
	)

	// Auto-migrate
	if err := db.AutoMigrate(&model.User{}); err != nil {
		panic("Auto migrate fail: " + err.Error())
	}
	log.Println("Auto Migrate success")

	// --- Dependency Injection ---
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, passwordHasher, jwtManager, mailer)
	authHandler := handler.NewAuthHandler(authService)

	// Setup Echo
	e := echo.New()
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Validator = &handler.CustomValidator{Validator: validator.New()}
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	// Routing
	api := e.Group("/api/auth")
	{
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)
	}

	// Jalankan Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Auth service listening on port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}