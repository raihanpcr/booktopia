package model

import (
	"time"
	// Import package BSON dari driver MongoDB
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	// Gunakan primitive.ObjectID untuk tipe ID di MongoDB
	// Tag bson:"_id,omitempty" berarti field ini akan dipetakan ke _id di MongoDB
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title          string             `json:"title" bson:"title"`
	Author         string             `json:"author" bson:"author"`
	Publisher      string             `json:"publisher" bson:"publisher"`
	YearPublished  int                `json:"year_published" bson:"year_published"`
	Category       string             `json:"category" bson:"category"`
	Price          float64            `json:"price" bson:"price"`
	Status         string             `json:"status" bson:"status"`
	IsDonationOnly bool               `json:"is_donation_only" bson:"is_donation_only"`
	Description    string             `json:"description" bson:"description"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
}