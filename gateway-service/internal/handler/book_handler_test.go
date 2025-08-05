package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// Skenario 1: Tes GetBooks jika proxy ke backend berhasil
func TestGetBooks_ProxySuccess(t *testing.T) {
	// --- Arrange ---
	// 1. Buat server backend palsu yang merespons dengan sukses
	mockBackend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"title":"Buku Tes"}]`))
	}))
	defer mockBackend.Close()

	// 2. Siapkan Echo context
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/books", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// 3. Buat handler dengan URL menunjuk ke server backend palsu
	h := NewBookHandler(mockBackend.URL)

	// --- Act ---
	// Panggil fungsi handler yang ingin dites
	err := h.GetBooks(c)

	// --- Assert ---
	// Pastikan hasilnya sesuai harapan
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `[{"title":"Buku Tes"}]`, rec.Body.String())
}

// Skenario 2: Tes CreateBook jika proxy ke backend berhasil
func TestCreateBook_ProxySuccess(t *testing.T) {
	// --- Arrange ---
	// 1. Buat server backend palsu
	mockBackend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id":"123","title":"Buku Baru"}`))
	}))
	defer mockBackend.Close()

	// 2. Siapkan Echo context untuk request POST
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/books", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// 3. Buat handler
	h := NewBookHandler(mockBackend.URL)

	// --- Act ---
	err := h.CreateBook(c)

	// --- Assert ---
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

// Skenario 3: Tes jika backend tidak bisa dihubungi
func TestProxy_BackendError(t *testing.T) {
	// --- Arrange ---
	// 1. Buat server backend palsu lalu langsung matikan untuk simulasi error
	mockBackend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	mockBackend.Close()

	// 2. Siapkan Echo context
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/books", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// 3. Buat handler
	h := NewBookHandler(mockBackend.URL)

	// --- Act ---
	// Panggil salah satu fungsi proxy (contoh: GetBooks)
	h.GetBooks(c)

	// --- Assert ---
	// Pastikan gateway mengembalikan error yang benar
	assert.Equal(t, http.StatusBadGateway, rec.Code)
	assert.Contains(t, rec.Body.String(), "failed to reach book service")
}
