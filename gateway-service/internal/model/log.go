package model

import "time"

type ServiceLog struct {
	ID           uint      `gorm:"primaryKey"`
	Timestamp    time.Time `gorm:"not null"`
	ServiceName  string
	Method       string
	URI          string
	StatusCode   int
	LatencyMs    int64
	RemoteIP     string
	UserAgent    string
	ErrorMessage string
	RequestBody  string `json:"request_body"`
}