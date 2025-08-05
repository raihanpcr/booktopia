package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// BookDTO adalah representasi data buku dari REST API book-service.
// Ini adalah DTO yang dilihat oleh transaction-service.
type BookDTO struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Price  float64 `json:"price"`
	Status string  `json:"status"`
}

type BookServiceTemplateResponse struct {
	Data BookDTO `json:"data"`
}

// BookServiceClient adalah interface untuk klien HTTP ke book-service.
type BookServiceClient interface {
	GetBookByID(ctx context.Context, bookID string) (*BookDTO, error)
}

type bookServiceClient struct {
	client  *http.Client
	baseURL string // Contoh: http://localhost:8081
}

// NewBookServiceClient membuat instance baru dari HTTP client.
func NewBookServiceClient(baseURL string) BookServiceClient {
	return &bookServiceClient{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		baseURL: baseURL,
	}
}

// GetBookByID melakukan panggilan GET ke /api/books/:id
func (c *bookServiceClient) GetBookByID(ctx context.Context, bookID string) (*BookDTO, error) {
	url := fmt.Sprintf("%s/books/%s", c.baseURL, bookID)
	fmt.Println(url)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("book-service returned non-200 status: %d", resp.StatusCode)
	}

	// Decode ke struct template
	var serviceResponse BookServiceTemplateResponse
	if err := json.NewDecoder(resp.Body).Decode(&serviceResponse); err != nil {
		return nil, err
	}

	// Kembalikan hanya bagian datanya
	return &serviceResponse.Data, nil
}
