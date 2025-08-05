package main

import (
	"gateway-service/internal/model"
	"gateway-service/internal/repository"
	"log"
	"os"

	_ "gateway-service/docs"
	"gateway-service/internal/handler"
	customMiddleware "gateway-service/internal/middleware"
	route "gateway-service/routes"
	gifting_pb "gifting-service/proto"
	transaction_pb "transaction-service/proto"
	wallet_pb "wallet-service/proto"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title           Booktopia Gateway API
// @version         1.0
// @description     Gateway service yang mengelola dan meneruskan permintaan ke berbagai microservice di ekosistem Booktopia, termasuk auth-service, book-service, gifting-service, transaction-service, dan wallet-service.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Booktopia Dev Team
// @contact.email  dev@booktopia.local

// @host      34.101.226.106:8000
// @BasePath  /api

// @schemes    http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Memuat variabel konfigurasi dari file .env
	godotenv.Load()

	// Ambil semua alamat service tujuan dari .env
	port := os.Getenv("PORT")
	dbURL := os.Getenv("DATABASE_URL")
	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	bookServiceURL := os.Getenv("BOOK_SERVICE_URL")
	transactionServiceURL := os.Getenv("TRANSACTION_SERVICE_URL")
	walletServiceURL := os.Getenv("WALLET_SERVICE_URL")
	giftingServiceURL := os.Getenv("GIFTING_SERVICE_URL")

	if port == "" {
		port = "8000"
	}
	if dbURL == "" {
		log.Fatal("DATABASE_URL for logging is not set")
	}

	// === Koneksi Database GORM untuk Logging ===
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to logging database: %v", err)
	}
	// AutoMigrate untuk tabel log
	db.AutoMigrate(&model.ServiceLog{})
	// Buat instance log repository
	logRepo := repository.NewGormLogRepository(db)

	// === Buat koneksi KLIEN ke semua service backend ===

	// Klien untuk transaction-service (gRPC)
	transactionConn, err := grpc.Dial(transactionServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect to transaction-service: %v", err)
	}
	defer transactionConn.Close()
	transactionClient := transaction_pb.NewTransactionServiceClient(transactionConn)

	// Klien untuk wallet-service (gRPC)
	walletConn, err := grpc.Dial(walletServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect to wallet-service: %v", err)
	}
	defer walletConn.Close()
	walletClient := wallet_pb.NewWalletServiceClient(walletConn)

	// Klien untuk gifting-service (gRPC)
	giftingConn, err := grpc.Dial(giftingServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect to gifting-service: %v", err)
	}
	defer giftingConn.Close()
	giftingClient := gifting_pb.NewGiftingServiceClient(giftingConn)

	// === Setup Server Echo (API Gateway) ===
	e := echo.New()

	// Gunakan middleware Logrus yang sudah diinisialisasi dengan logRepo
	e.Use(customMiddleware.LogrusMiddleware(logRepo))
	e.Use(middleware.CORS()) // Middleware untuk mengizinkan Cross-Origin Resource Sharing

	// Inisialisasi semua handler, menyuntikkan (inject) URL atau gRPC client yang dibutuhkan
	authHandler := handler.NewAuthHandler(authServiceURL)
	bookHandler := handler.NewBookHandler(bookServiceURL)
	transactionHandler := handler.NewTransactionHandler(transactionClient)
	walletHandler := handler.NewWalletHandler(walletClient)
	giftingHandler := handler.NewGiftingHandler(giftingClient)

	// Mendaftarkan semua route API dari file terpisah
	route.SetupRoutes(e, authHandler, bookHandler, transactionHandler, walletHandler, giftingHandler)

	// Mendaftarkan route untuk halaman dokumentasi Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Memulai server pada port yang telah ditentukan
	log.Printf("API Gateway listening on port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
