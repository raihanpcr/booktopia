package app

import "errors"

// Definisi Error kustom untuk lapisan aplikasi Auth.
// dto error message
var (
	ErrValidation   = errors.New("validasi gagal: satu atau lebih field kosong")
	ErrEmailExist   = errors.New("email sudah terdaftar")
	ErrUnauthorized = errors.New("email atau password salah")
)