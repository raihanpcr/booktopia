package dto

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