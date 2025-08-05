package dto

import "book-service/internal/model"

// ToBookModel mengubah DTO CreateBookRequest menjadi model internal.
func (r *CreateBookRequest) ToBookModel() *model.Book {
	return &model.Book{
		Title:          r.Title,
		Author:         r.Author,
		Publisher:      r.Publisher,
		YearPublished:  r.YearPublished,
		Category:       r.Category,
		Price:          r.Price,
		IsDonationOnly: r.IsDonationOnly,
		Description:    r.Description,
	}
}

// ToBookModel mengubah DTO UpdateBookRequest menjadi model internal.
func (r *UpdateBookRequest) ToBookModel() *model.Book {
	return &model.Book{
		Title:          r.Title,
		Author:         r.Author,
		Publisher:      r.Publisher,
		YearPublished:  r.YearPublished,
		Category:       r.Category,
		Price:          r.Price,
		Status:         r.Status,
		IsDonationOnly: r.IsDonationOnly,
		Description:    r.Description,
	}
}

// ToBookResponse mengubah model internal menjadi DTO response.
func ToBookResponse(book model.Book) BookResponse {
	return BookResponse{
		ID:             book.ID.Hex(), // Ubah ObjectID ke string
		Title:          book.Title,
		Author:         book.Author,
		Publisher:      book.Publisher,
		YearPublished:  book.YearPublished,
		Category:       book.Category,
		Price:          book.Price,
		Status:         book.Status,
		IsDonationOnly: book.IsDonationOnly,
		Description:    book.Description,
		CreatedAt:      book.CreatedAt,
	}
}

// ToBookResponseList mengubah slice model menjadi slice DTO response.
func ToBookResponseList(books []model.Book) []BookResponse {
	var bookResponses []BookResponse
	for _, b := range books {
		bookResponses = append(bookResponses, ToBookResponse(b))
	}
	return bookResponses
}