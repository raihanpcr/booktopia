package handler

import (
	"auth-service/internal/auth/app"
	"auth-service/internal/dto"
	"auth-service/internal/service"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// CustomValidator untuk mengintegrasikan validator dengan Echo.
type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

// AuthHandler menangani request HTTP untuk otentikasi.
type AuthHandler struct {
	Service service.AuthService
}

// NewAuthHandler adalah constructor untuk AuthHandler.
func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{Service: service}
}

// Register godoc
// @Summary      Register user
// @Description  Membuat akun pengguna baru
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RegisterRequest true "Register Request"
// @Success      201  {object}  dto.TemplateRegisterResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      409  {object}  dto.ErrorResponse
// @Router       /register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	// 1. Bind & Validate request body
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message: "Invalid Request",
			Error: err.Error(),
		})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message: "Invalid Request",
			Error: err.Error(),
		})
	}

	// 2. Panggil service untuk menjalankan logika bisnis
	resp, err := h.Service.Register(c.Request().Context(), req)
	if err != nil {
		// 3. Terjemahkan error dari service ke response HTTP
		return errToHTTP(c, err)
	}

	return c.JSON(http.StatusCreated, dto.TemplateRegisterResponse{
		StatusCode: http.StatusCreated,
		Message: "Registration successful",
		Data: *resp,
	})
}

// Login godoc
// @Summary      Login user
// @Description  Login pengguna dan mengembalikan token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "Login Request"
// @Success      200  {object}  dto.TemplateLoginResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Router       /login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	// 1. Bind & Validate request body
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message: "Invalid Request",
			Error: err.Error(),
		})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:"Invalid Request",
			Error: err.Error(),
		})
	}

	// 2. Panggil service untuk menjalankan logika bisnis
	resp, err := h.Service.Login(c.Request().Context(), req)
	if err != nil {
		// 3. Terjemahkan error dari service ke response HTTP
		return errToHTTP(c, err)
	}
	return c.JSON(http.StatusOK, dto.TemplateLoginResponse{
		StatusCode: http.StatusOK,
		Message: "Login Successful",
		Data: *resp,
	})
}

// errToHTTP adalah fungsi helper untuk memetakan error dari service
// ke response HTTP dengan status code yang sesuai.
func errToHTTP(c echo.Context, err error) error {
	switch err {
	case app.ErrEmailExist:
		return c.JSON(http.StatusConflict, dto.ErrorResponse{
			StatusCode: http.StatusConflict,
			Message: "Email already exist",
			Error: err.Error(),
		})
	case app.ErrUnauthorized:
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message: "Incorrect email or password",
			Error: err.Error(),
		})
	default:
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message: "Internal error",
			Error: err.Error(),})
	}
}