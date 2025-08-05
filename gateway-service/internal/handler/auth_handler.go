package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authServiceURL string
}

func NewAuthHandler(authServiceURL string) *AuthHandler {
	return &AuthHandler{authServiceURL: authServiceURL}
}

// @Summary Register User
// @Description Meneruskan permintaan register ke Auth Service
// @Tags Gateway - Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Data pendaftaran user"
// @Success 201 {object} dto.TemplateRegisterResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      409  {object}  dto.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	// Panggil fungsi proxy yang sama untuk menghindari duplikasi kode
	return h.proxyToAuth(c)
}

// @Summary Login User
// @Description Meneruskan permintaan login ke Auth Service
// @Tags Gateway - Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Data login user"
// @Success 200 {object} dto.TemplateLoginResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      409  {object}  dto.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	// Panggil fungsi proxy yang sama untuk menghindari duplikasi kode
	return h.proxyToAuth(c)
}

// proxyToAuth adalah fungsi internal untuk logika proxy
func (h *AuthHandler) proxyToAuth(c echo.Context) error {
	targetURL, _ := url.Parse(fmt.Sprintf("%s%s", h.authServiceURL, c.Request().URL.Path))

	proxyReq, err := http.NewRequest(c.Request().Method, targetURL.String(), c.Request().Body)
	if err != nil {
		return err
	}
	proxyReq.Header = make(http.Header)
	for k, v := range c.Request().Header {
		if k == "Authorization" {
			continue
		}
		proxyReq.Header[k] = v
	}

	client := http.DefaultClient
	resp, err := client.Do(proxyReq)
	if err != nil {
		return c.JSON(http.StatusBadGateway, map[string]string{
			"error": "failed to reach auth service",
			"detail": err.Error(),
		})
	}
	defer resp.Body.Close()

	c.Response().WriteHeader(resp.StatusCode)
	io.Copy(c.Response().Writer, resp.Body)

	return nil
}
