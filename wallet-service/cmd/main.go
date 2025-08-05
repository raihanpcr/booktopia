package main

import (
	"context"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"wallet-service/internal/model"
	"wallet-service/internal/repository"
	"wallet-service/internal/server"
	"wallet-service/internal/service"
	"wallet-service/internal/worker"
	pb "wallet-service/proto"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Ambil konfigurasi dari environment
	dbURL := os.Getenv("DATABASE_URL")
	grpcPort := os.Getenv("GRPC_PORT")
	kafkaURL := os.Getenv("KAFKA_URL")

	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}
	if kafkaURL == "" {
		log.Fatal("KAFKA_URL is not set")
	}
	if grpcPort == "" {
		grpcPort = "50053"
	}

	// Koneksi ke PostgreSQL menggunakan GORM dengan DATABASE_URL
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// AutoMigrate untuk membuat tabel
	log.Println("Running migrations...")
	db.AutoMigrate(&model.User{}, &model.TopUp{})

	// Inisialisasi dependensi
	repo := repository.NewGormRepository(db)
	svc := service.NewWalletService(repo)
	grpcServer := server.NewGrpcServer(svc)

	// === Jalankan Kafka Consumer ===
	// Gunakan context untuk bisa mematikan consumer saat aplikasi berhenti
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	worker.StartConsumer(ctx, kafkaURL, "transaction_created", svc)

	// Setup dan jalankan server gRPC
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", grpcPort, err)
	}

	s := grpc.NewServer()
	pb.RegisterWalletServiceServer(s, grpcServer)

	log.Printf("Wallet gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
