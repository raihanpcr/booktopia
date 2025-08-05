package dto

import "time"

// CreateBookRequest adalah DTO untuk membuat buku baru.
// Tidak ada ID, Status, atau CreatedAt karena itu diatur oleh server.
type CreateBookRequest struct {
	Title          string  `json:"title" validate:"required"`
	Author         string  `json:"author" validate:"required"`
	Publisher      string  `json:"publisher"`
	YearPublished  int     `json:"year_published"`
	Category       string  `json:"category"`
	Price          float64 `json:"price" validate:"gte=0"`
	IsDonationOnly bool    `json:"is_donation_only"`
	Description    string  `json:"description"`
}

// UpdateBookRequest adalah DTO untuk memperbarui buku.
// Mirip dengan Create, tapi semua field bisa jadi opsional tergantung logika bisnis.
type UpdateBookRequest struct {
	Title          string  `json:"title" validate:"required"`
	Author         string  `json:"author" validate:"required"`
	Publisher      string  `json:"publisher"`
	YearPublished  int     `json:"year_published"`
	Category       string  `json:"category"`
	Price          float64 `json:"price" validate:"gte=0"`
	Status         string  `json:"status" validate:"oneof=available unavailable"` // Validasi status
	IsDonationOnly bool    `json:"is_donation_only"`
	Description    string  `json:"description"`
}

// BookResponse adalah DTO untuk data buku yang dikirim ke klien.
// ID di sini adalah string agar mudah dikonsumsi oleh JSON.
type BookResponse struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	Author         string    `json:"author"`
	Publisher      string    `json:"publisher"`
	YearPublished  int       `json:"year_published"`
	Category       string    `json:"category"`
	Price          float64   `json:"price"`
	Status         string    `json:"status"`
	IsDonationOnly bool      `json:"is_donation_only"`
	Description    string    `json:"description"`
	CreatedAt      time.Time `json:"created_at"`
}

type DeleteResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message"`
}

type BookCreateResponse struct {
	StatusCode int              `json:"status_code" validate:"required" example:"201"`
	Message    string           `json:"message" validate:"required" example:"Create user success"`
	Data       BookResponse 	`json:"data"`
}

type BookGetResponse struct {
	StatusCode int              `json:"status_code" validate:"required" example:"201"`
	Message    string           `json:"message" validate:"required" example:"Create user success"`
	Data       []BookResponse `json:"data"`
}