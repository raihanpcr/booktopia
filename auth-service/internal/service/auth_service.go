package service

import (
	"context"
	"fmt"

	"auth-service/internal/auth/app"
	"auth-service/internal/dto"
	"auth-service/internal/model"
	"auth-service/internal/repository"
	"auth-service/pkg/mail"
)

// AuthService adalah interface untuk logika bisnis.
type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
}

// authService adalah implementasi dari AuthService.
type authService struct {
	repo       repository.UserRepository
	hasher     app.PasswordHasher
	jwtManager app.JWTManager
	mailer    	mail.Mailer
}

// NewAuthService adalah constructor untuk membuat instance authService.
func NewAuthService(
	r repository.UserRepository,
	h app.PasswordHasher,
	j app.JWTManager,
	m mail.Mailer,
) AuthService {
	return &authService{repo: r, hasher: h, jwtManager: j, mailer: m}
}

// Register menangani logika registrasi, hashing, penyimpanan, dan pengiriman email.
func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error) {
	// Cek apakah email sudah terdaftar
	_, err := s.repo.GetByEmail(req.Email)
	if err == nil {
		return nil, app.ErrEmailExist
	}

	// Hash password
	hashedPassword, err := s.hasher.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	// Buat entitas user
	user := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "pembeli",
	}

	// Simpan user ke database
	createdUser, err := s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	// Kirim email selamat datang (secara non-blocking)
	// err = s.mailer.SendWelcomeEmail(createdUser.Email, createdUser.Name)
	go s.mailer.SendWelcomeEmail(createdUser.Email, createdUser.Name)

	// Kembalikan response DTO
	return &dto.RegisterResponse{
		ID:    createdUser.ID,
		Name:  createdUser.Name,
		Email: createdUser.Email,
	}, nil
}

// Login menangani logika verifikasi kredensial dan pembuatan token JWT.
func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	// Cari user berdasarkan email
	user, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		return nil, app.ErrUnauthorized
	}

	// Verifikasi password
	if !s.hasher.Check(req.Password, user.Password) {
		return nil, app.ErrUnauthorized
	}

	// Generate token JWT
	token, err := s.jwtManager.GenerateToken(fmt.Sprintf("%d", user.ID), user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	// Kembalikan response DTO
	return &dto.AuthResponse{
		Email: user.Email,
		Role: user.Role,
		Token: token,
	}, nil
}