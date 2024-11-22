package repository

import (
	"context"

	"test-backend-altech/model/domain"
	"test-backend-altech/model/web/response"
	"test-backend-altech/repository/query"

	"github.com/jackc/pgx/v5"
)

type bookRepository struct {
	db        Store
	BookQuery query.BookQuery
}

type BookRepository interface {
	CreateBook(c context.Context, book domain.Book) error
	UpdateBook(c context.Context, id string, book domain.UpdateBook) error
	FindByID(c context.Context, id string) (response.BookResponse, error)
	ValidateBookTitle(c context.Context, name string) (domain.ValidateBookTitle, error)
	FindAllBook(c context.Context) ([]response.BookResponse, error)
	DeleteBook(c context.Context, id string) error
}

func NewBookRepository(db Store, q query.BookQuery) BookRepository {
	return &bookRepository{
		db:        db,
		BookQuery: q,
	}
}

func (r *bookRepository) CreateBook(c context.Context, book domain.Book) error {
	var err error

	err = r.db.WithTransaction(c, func(tx pgx.Tx) error {
		if err = r.BookQuery.CreateBook(c, tx, book); err != nil {
			return err
		}

		return nil
	})

	return err
}

func (r *bookRepository) UpdateBook(c context.Context, id string, book domain.UpdateBook) error {
	var err error

	err = r.db.WithTransaction(c, func(tx pgx.Tx) error {
		if err = r.BookQuery.UpdateBook(c, tx, id, book); err != nil {
			return err
		}
		return nil
	})

	return err
}

func (r *bookRepository) FindByID(c context.Context, id string) (response.BookResponse, error) {
	var err error
	var book response.BookResponse

	err = r.db.WithTransaction(c, func(tx pgx.Tx) error {
		book, err = r.BookQuery.FindByID(c, tx, id)
		return err
	})

	return book, err
}

func (r *bookRepository) ValidateBookTitle(c context.Context, name string) (domain.ValidateBookTitle, error) {
	var err error
	var book domain.ValidateBookTitle

	err = r.db.WithTransaction(c, func(tx pgx.Tx) error {
		book, err = r.BookQuery.ValidateBookTitle(c, tx, name)
		return err
	})

	return book, err
}

func (r *bookRepository) FindAllBook(c context.Context) ([]response.BookResponse, error) {
	var err error
	var books []response.BookResponse

	err = r.db.WithTransaction(c, func(tx pgx.Tx) error {
		books, err = r.BookQuery.FindAllBook(c, tx)
		return err
	})

	return books, err
}

func (r *bookRepository) DeleteBook(c context.Context, id string) error {
	var err error

	err = r.db.WithTransaction(c, func(tx pgx.Tx) error {
		err = r.BookQuery.DeleteBook(c, tx, id)
		return err
	})

	return err
}
