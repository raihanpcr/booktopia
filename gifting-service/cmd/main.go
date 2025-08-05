package main

import (
	"context"
	"gifting-service/internal/model"
	"gifting-service/internal/repository"
	"gifting-service/internal/server"
	"gifting-service/internal/service"
	"gifting-service/pkg/client"
	pb "gifting-service/proto"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	grpcPort := os.Getenv("GRPC_PORT")
	bookServiceURL := os.Getenv("BOOK_SERVICE_URL")

	if grpcPort == "" {
		grpcPort = "50054"
	}

	// Koneksi Database
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	db.AutoMigrate(&model.EbookGiftLog{})

	// Inisialisasi Dependensi
	bookClient := client.NewBookServiceClient(bookServiceURL)
	repo := repository.NewGormRepository(db)
	svc := service.NewGiftingService(repo, bookClient)
	grpcServer := server.NewGrpcServer(svc)

	// Setup dan Jalankan Scheduller
	c := cron.New()
	// Jalankan tugas setiap hari tengah malam
	c.AddFunc("0 0 * * *", func() {
		log.Println("Running scheduled task: Expiring old gifts...")
		svc.ExpiredOldGifts(context.Background())
	})
	c.Start()
	log.Println("Cron Scheduller has been started.")

	// Jalankan Server gRPC
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGiftingServiceServer(s, grpcServer)

	log.Printf("Gifting gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
