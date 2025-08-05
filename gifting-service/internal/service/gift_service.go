package service

import (
	"context"
	"errors"
	"gifting-service/internal/model"
	"gifting-service/internal/repository"
	"gifting-service/pkg/client" // HTTP client ke book-service
	pb "gifting-service/proto"
	"log"
	"strconv"
)

type GiftingService interface {
	SendGift(ctx context.Context, req *pb.SendGiftRequest) (*model.EbookGiftLog, error)
	ExpiredOldGifts(ctx context.Context)
}

type giftingService struct {
	repo       repository.GiftingRepository
	bookClient client.BookServiceClient
}

func NewGiftingService(repo repository.GiftingRepository, bookClient client.BookServiceClient) GiftingService {
	return &giftingService{repo: repo, bookClient: bookClient}
}

func (s *giftingService) SendGift(ctx context.Context, req *pb.SendGiftRequest) (*model.EbookGiftLog, error) {
	// 1. Validasi buku ke book-service
	book, err := s.bookClient.GetBookByID(ctx, req.BookId)
	if err != nil {
		return nil, errors.New("book not found")
	}
	// Pastikan buku tersebut memang untuk donasi
	if book.IsDonationOnly {
		return nil, errors.New("this book cannot be gifted")
	}

	donorID, _ := strconv.ParseUint(req.DonorId, 10, 32)
	bookID, _ := strconv.ParseUint(req.BookId, 10, 32)

	// 2. Buat entitas hadiah
	gift := &model.EbookGiftLog{
		DonorID:        uint(donorID),
		RecipientEmail: req.RecipientEmail,
		BookID:         uint(bookID),
		Message:        req.Message,
		Status:         "pending",
	}

	// 3. Simpan ke database
	return s.repo.CreateGift(ctx, gift)
}

// Scheduller (CronJob)
func (s *giftingService) ExpiredOldGifts(ctx context.Context) {
	log.Println("Scheduler running: Expiring old gifts...")

	// Set hadiah untuk expired setelah 7 hari
	rowsAffected, err := s.repo.ExpiredOldGifts(ctx, 7)
	if err != nil {
		log.Printf("Error expiring old gifts: %v", err)
		return
	}

	log.Printf("Scheduler finished: %d old gifts expired.", rowsAffected)
}
