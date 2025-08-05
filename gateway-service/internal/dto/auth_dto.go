package dto

import "time"

// RegisterRequest merepresentasikan struktur body permintaan untuk pendaftaran pengguna baru.
// Struct ini digunakan untuk binding data JSON dari permintaan HTTP.
// Tag `json` digunakan untuk mapping field Go ke nama kunci JSON.
// Tag `validate` digunakan oleh library validator (misalnya, go-playground/validator)
// untuk menerapkan aturan validasi pada input yang diterima.
type RegisterRequest struct {
	Name     string `json:"name" validate:"required" example:"John Doe"`                    // Nama lengkap pengguna. `required` berarti tidak boleh kosong.
	Email    string `json:"email" validate:"required,email" example:"john.doe@example.com"` // Alamat email pengguna. `required` dan `email` (format email valid).
	Password string `json:"password" validate:"required,min=8" example:"password123"`       // Kata sandi pengguna. `required` dan `min=8` (minimal 8 karakter).
}

type RegisterResponse struct {
	ID    uint   `json:"id" validate:"required" example:"1"`
	Name  string `json:"name" validate:"required" example:"John Doe"`
	Email string `json:"email" validate:"required,email" example:"john.doe@example.com"`
}

// LoginRequest merepresentasikan struktur body permintaan untuk login pengguna.
// Struct ini digunakan untuk binding data JSON dari permintaan HTTP.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"john.doe@example.com"` // Alamat email pengguna. `required` dan `email`.
	Password string `json:"password" validate:"required" example:"password123"`             // Kata sandi pengguna. `required`.
}

// AuthResponse merepresentasikan struktur respons untuk operasi otentikasi (Register atau Login).
// Field-field menggunakan `omitempty` agar tidak disertakan dalam respons JSON jika nilainya kosong.
type AuthResponse struct {
	ID    uint `json:"id,omitempty" example:"1"` // ID pengguna, hanya disertakan jika ada (saat Register).
	Name  string `json:"name,omitempty" example:"John Doe"`               // Nama pengguna, hanya disertakan jika ada (saat Register).
	Email string `json:"email,omitempty" example:"john.doe@example.com"`  // Alamat email pengguna.
	Token string `json:"token,omitempty" example:"eyJhbGciOiJIUzI1Ni..."` // Token JWT, hanya disertakan jika ada (saat Login atau Register jika langsung login).
}

// ErrorResponse merepresentasikan struktur respons standar untuk error API.
// Ini membantu dalam memberikan format error yang konsisten kepada klien.
type ErrorResponse struct {
	StatusCode int `json:"status_code" validate:"required" example:"401"`
	Message string `json:"message" validate:"required" example:"Internal Server Error"`
	Error string `json:"error" example:"Pesan kesalahan yang deskriptif"`
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Saldo     float64   `json:"saldo"`
	CreatedAt time.Time`json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetAllUsersResponse struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Data    []UserResponse `json:"data"`
}

type TemplateLoginResponse struct {
	StatusCode int              `json:"status_code" validate:"required" example:"201"`
	Message    string           `json:"message" validate:"required" example:"Create user success"`
	Data       AuthResponse 	`json:"data"`
}

type TemplateRegisterResponse struct {
	StatusCode int              `json:"status_code" validate:"required" example:"201"`
	Message    string           `json:"message" validate:"required" example:"Create user success"`
	Data       RegisterResponse `json:"data"`
}