package dto

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
	Role string `json:"role,omitempty" example:"pembeli"` 
	Token string `json:"token,omitempty" example:"eyJhbGciOiJIUzI1Ni..."` // Token JWT, hanya disertakan jika ada (saat Login atau Register jika langsung login).
}

type TemplateLoginResponse struct {
	StatusCode int              `json:"status_code" validate:"required" example:"201"`
	Message    string           `json:"message" validate:"required" example:"Create user success"`
	Data       AuthResponse 	`json:"data"`
}