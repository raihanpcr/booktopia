package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5" // Menggunakan versi 5 dari library JWT
)

// Claims adalah struktur kustom yang mendefinisikan klaim (payload) yang akan
// disertakan dalam token JWT. Ini mengimplementasikan antarmuka jwt.Claims
// dengan menyematkan jwt.RegisteredClaims, yang menyediakan klaim standar JWT.
type Claims struct {
	UserID               string `json:"user_id"`
	Email                string `json:"email"`
	Role                 string `json:"role"`
	jwt.RegisteredClaims        // Klaim standar JWT seperti Issuer, Subject, Audience, ExpiresAt, dll.
}

// Manager adalah struktur yang bertanggung jawab untuk membuat dan memverifikasi token JWT.
// Ini menyimpan secret key yang digunakan untuk menandatangani dan memverifikasi token.
type Manager struct {
	Secret string // Secret key yang digunakan untuk menandatangani dan memverifikasi JWT
}

// NewManager adalah konstruktor untuk Manager.
// Fungsi ini menginisialisasi Manager dengan secret key yang diberikan.
func NewManager(secret string) *Manager {
	return &Manager{Secret: secret}
}

// GenerateToken membuat dan mengembalikan string token JWT yang baru.
// Token ini akan berisi userID dan email sebagai klaim kustom,
// serta klaim standar seperti waktu kedaluwarsa (exp).
func (j *Manager) GenerateToken(userID, email, role string) (string, error) {
	// Tentukan waktu kedaluwarsa token (misalnya, 24 jam dari sekarang)
	expirationTime := time.Now().Add(time.Hour * 24)

	// Buat objek Claims kustom
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime), // Waktu kedaluwarsa token
			IssuedAt:  jwt.NewNumericDate(time.Now()),     // Waktu token diterbitkan
			// Anda bisa menambahkan klaim standar lainnya di sini jika diperlukan
			// Issuer: "auth-service",
			// Subject: userID,
			// Audience: []string{"web", "mobile"},
		},
	}

	// Buat token JWT baru dengan metode penandatanganan HS256 dan klaim yang sudah dibuat
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Tandatangani token menggunakan secret key dan kembalikan string token
	signedToken, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return "", fmt.Errorf("gagal menandatangani token: %w", err)
	}
	return signedToken, nil
}

// Verify memverifikasi string token JWT.
// Fungsi ini memeriksa validitas tanda tangan token dan waktu kedaluwarsa.
// Jika token valid, ia mengembalikan pointer ke objek Claims yang berisi payload token.
// Jika token tidak valid atau ada masalah, ia mengembalikan error.
func (j *Manager) Verify(tokenString string) (*Claims, error) {
	// Parse token string menggunakan secret key dan fungsi callback untuk validasi
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Memastikan metode penandatanganan adalah HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metode penandatanganan tidak valid: %v", token.Header["alg"])
		}
		// Mengembalikan secret key sebagai slice of byte
		return []byte(j.Secret), nil
	})

	if err != nil {
		// Menangani berbagai jenis error JWT
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token kadaluarsa")
		}
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, fmt.Errorf("tanda tangan token tidak valid")
		}
		return nil, fmt.Errorf("token tidak valid: %w", err)
	}

	// Memeriksa apakah token secara keseluruhan valid
	if !token.Valid {
		return nil, fmt.Errorf("token tidak valid")
	}

	// Mendapatkan klaim dari token yang sudah diverifikasi
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("tipe klaim tidak valid")
	}

	return claims, nil
}
