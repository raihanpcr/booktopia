package dto

import "time"

// RegisterRequest merepresentasikan struktur body permintaan untuk pendaftaran pengguna baru.
// Struct ini digunakan untuk binding data JSON dari permintaan HTTP.
// Tag `json` digunakan untuk mapping field Go ke nama kunci JSON.
// Tag `validate` digunakan oleh library validator (misalnya, go-playground/validator)
// untuk menerapkan aturan validasi pada input yang diterima.









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

