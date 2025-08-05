package dto

// ErrorResponse adalah struktur standar untuk pesan error API.
type ErrorResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"` // Tipe 'any' (atau interface{}) agar fleksibel
}