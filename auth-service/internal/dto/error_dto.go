package dto

// ErrorResponse merepresentasikan struktur respons standar untuk error API.
// Ini membantu dalam memberikan format error yang konsisten kepada klien.
type ErrorResponse struct {
	StatusCode int `json:"status_code" validate:"required" example:"401"`
	Message string `json:"message" validate:"required" example:"401"`
	Error string `json:"error" example:"Pesan kesalahan yang deskriptif"` // Pesan error yang menjelaskan masalah.
}