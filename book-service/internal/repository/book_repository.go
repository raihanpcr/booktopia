package repository

import (
	"context"
	"book-service/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// BookRepository mendefinisikan kontrak untuk interaksi database
type BookRepository interface {
	Create(ctx context.Context, book *model.Book) error
	FindAll(ctx context.Context) ([]model.Book, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*model.Book, error)
	Update(ctx context.Context, book *model.Book) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type bookRepository struct {
	collection *mongo.Collection
}

func NewBookRepository(collection *mongo.Collection) BookRepository {
	return &bookRepository{collection: collection}
}

// Create menyimpan satu buku baru
func (r *bookRepository) Create(ctx context.Context, book *model.Book) error {
	_, err := r.collection.InsertOne(ctx, book)
	return err
}

// FindAll mengambil semua buku yang tersedia
func (r *bookRepository) FindAll(ctx context.Context) ([]model.Book, error) {
	filter := bson.M{"status": "available"}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var books []model.Book
	if err = cursor.All(ctx, &books); err != nil {
		return nil, err
	}
	return books, nil
}

// FindByID mencari satu buku berdasarkan ID
func (r *bookRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Book, error) {
	var book model.Book
	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&book)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Dokumen tidak ditemukan adalah kasus yang valid
		}
		return nil, err
	}
	return &book, nil
}

// Update memperbarui dokumen buku yang ada
func (r *bookRepository) Update(ctx context.Context, book *model.Book) error {
	filter := bson.M{"_id": book.ID}
	// bson.M{"$set": book} akan memperbarui semua field di dokumen
	update := bson.M{"$set": book}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete menghapus dokumen buku berdasarkan ID
func (r *bookRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}