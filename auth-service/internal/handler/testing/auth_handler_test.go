package testing_test

import (
	"auth-service/internal/auth/app"
	"auth-service/internal/dto"
	"auth-service/internal/handler"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RegisterResponse), args.Error(1)
}

func (m *MockAuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.AuthResponse), args.Error(1)
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func newTestContext(method, path string, body interface{}) (echo.Context, *httptest.ResponseRecorder){
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(method, path, bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	return e.NewContext(req, rec), rec
}

func TestAuthHandler_Register_Success(t *testing.T) {
	mockService := new(MockAuthService)
	handler := handler.NewAuthHandler(mockService)

	req := dto.RegisterRequest{
		Name:     "John",
		Email:    "john@example.com",
		Password: "secret123",
	}
	resp := dto.RegisterResponse{
		ID:    1,
		Name:  "John",
		Email: "john@example.com",
	}

	mockService.On("Register", mock.Anything, req).Return(&resp, nil)

	c, rec := newTestContext(http.MethodPost, "/register", req)
	err := handler.Register(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var got dto.TemplateRegisterResponse
	err = json.Unmarshal(rec.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, got.StatusCode)
	assert.Equal(t, "Registration successful", got.Message)
	assert.Equal(t, resp, got.Data)
}

func TestAuthHandler_Login_Success(t *testing.T) {
	mockService := new(MockAuthService)
	handler := handler.NewAuthHandler(mockService)

	req := dto.LoginRequest{
		Email:    "john@example.com",
		Password: "secret123",
	}
	resp := dto.AuthResponse{
		ID:    1,
		Name:  "John",
		Email: "john@example.com",
		Token: "jwt-token",
	}

	mockService.On("Login", mock.Anything, req).Return(&resp, nil)

	c, rec := newTestContext(http.MethodPost, "/login", req)
	err := handler.Login(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var got dto.TemplateLoginResponse
	err = json.Unmarshal(rec.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, got.StatusCode)
	assert.Equal(t, "Login Successful", got.Message)
	assert.Equal(t, resp, got.Data)
}


// Duplicate email
func TestAuthHandler_Register_EmailAlreadyExist(t *testing.T) {

	mockService := new(MockAuthService)
	handler := handler.NewAuthHandler(mockService)

	req := dto.RegisterRequest{
		Name:     "John",
		Email:    "john@example.com",
		Password: "secret123",
	}

	mockService.On("Register", mock.Anything, req).Return(nil, app.ErrEmailExist)

	c, rec := newTestContext(http.MethodPost, "/register", req)
	err := handler.Register(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusConflict, rec.Code)

	var got dto.ErrorResponse
	json.Unmarshal(rec.Body.Bytes(), &got)
	assert.Equal(t, "email sudah terdaftar", got.Error)
}

// Login Gagal
func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	mockService := new(MockAuthService)
	handler := handler.NewAuthHandler(mockService)

	req := dto.LoginRequest{
		Email:    "john@example.com",
		Password: "wrong-password",
	}

	mockService.On("Login", mock.Anything, req).Return(nil, app.ErrUnauthorized)

	c, rec := newTestContext(http.MethodPost, "/login", req)
	err := handler.Login(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var got dto.ErrorResponse
	json.Unmarshal(rec.Body.Bytes(), &got)
	assert.Equal(t, "email atau password salah", got.Error)
}
