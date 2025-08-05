package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres" // 1. Ganti driver
	"gorm.io/gorm"

	"transaction-service/internal/model" // 2. Import model untuk AutoMigrate
	"transaction-service/internal/repository"
	"transaction-service/internal/server"
	"transaction-service/internal/service"
	"transaction-service/pkg/client"
	"transaction-service/pkg/messagebroker"
	pb "transaction-service/proto"
	wallet_pb "wallet-service/proto"
)

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	grpcPort := os.Getenv("GRPC_PORT")
	bookServiceURL := os.Getenv("BOOK_SERVICE_URL")
	walletServiceURL := os.Getenv("WALLET_SERVICE_URL")
	kafkaURL := os.Getenv("KAFKA_URL")
	fmt.Println(kafkaURL)

	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	if bookServiceURL == "" {
		log.Fatal("BOOK_SERVICE_URL is not set")
	}
	if walletServiceURL == "" {
		log.Fatal("WALLET_SERVICE_URL is not set")
	}
	if kafkaURL == "" {
		log.Fatal("KAFKA_URL is not set")
	}
	if grpcPort == "" {
		grpcPort = "50052"
	}

	// 3. Koneksi ke PostgreSQL menggunakan GORM
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	// 4. Jalankan AutoMigrate
	log.Println("Running migrations for transaction service...")
	db.AutoMigrate(&model.Transaction{}, &model.TransactionDetail{})

	// Koneksi KLIEN ke wallet-service
	walletConn, err := grpc.Dial(walletServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect to wallet-service: %v", err)
	}
	defer walletConn.Close()
	walletClient := wallet_pb.NewWalletServiceClient(walletConn)

	// Inisialisasi dependensi
	bookClient := client.NewBookServiceClient(bookServiceURL)
	kafkaProducer := messagebroker.NewKafkaProducer(kafkaURL)

	// 5. Gunakan GORM repository
	repo := repository.NewGormRepository(db)
	svc := service.NewTransactionService(repo, bookClient, walletClient, kafkaProducer)
	grpcServer := server.NewGrpcServer(svc)

	// Setup dan jalankan server gRPC
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTransactionServiceServer(s, grpcServer)

	log.Printf("Transaction gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
