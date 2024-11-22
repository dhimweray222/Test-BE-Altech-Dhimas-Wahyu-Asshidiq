package domain

import (
	"test-backend-altech/model/web/response"
	"time"

	"github.com/google/uuid"
)

type Book struct {
	Id          string
	Title       string
	Description string
	PublishDate string
	AuthorId    string
	CreatedAt   time.Time
}

func (book *Book) GenerateID() {
	uuid := uuid.New().String()
	book.Id = uuid
}

func (book *Book) ToBookResponse() response.BookResponse {
	return response.BookResponse{
		Id:          book.Id,
		Title:       book.Title,
		Description: book.Description,
		PublishDate: book.PublishDate,
		AuthorId:    book.AuthorId,
	}
}

type ValidateBookTitle struct {
	Title string `json:"name"`
}

type UpdateBook struct {
	Title       string `json:"title"`
	PublishDate string `json:"publish_date"`
	Description string `json:"description"`
	AuthorId    string `json:"author_id"`
}
