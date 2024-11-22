package repository

import (
	"context"

	"test-backend-altech/model/domain"
	"test-backend-altech/repository/query"

	"github.com/jackc/pgx/v5"
)

type authorRepository struct {
	db          Store
	AuthorQuery query.AuthorQuery
}

type AuthorRepository interface {
	CreateAuthor(c context.Context, author domain.Author) error
	UpdateAuthor(c context.Context, id string, author domain.UpdateAuthor) error
	FindByID(c context.Context, id string) (domain.Author, error)
	ValidateAuthorName(c context.Context, name string) (domain.ValidateAuthorName, error)
	FindAllAuthor(c context.Context) ([]domain.Author, error)
	DeleteAuthor(c context.Context, id string) error
}

func NewAuthorRepository(db Store, q query.AuthorQuery) AuthorRepository {
	return &authorRepository{
		db:          db,
		AuthorQuery: q,
	}
}

func (r *authorRepository) CreateAuthor(c context.Context, author domain.Author) error {
	var err error

	err = r.db.WithTransaction(c, func(tx pgx.Tx) error {
		if err = r.AuthorQuery.CreateAuthor(c, tx, author); err != nil {
			return err
		}
		return nil
	})

	return err
}

func (r *authorRepository) UpdateAuthor(c context.Context, id string, author domain.UpdateAuthor) error {
	var err error

	err = r.db.WithTransaction(c, func(tx pgx.Tx) error {
		if err = r.AuthorQuery.UpdateAuthor(c, tx, id, author); err != nil {
			return err
		}
		return nil
	})

	return err
}

func (r *authorRepository) FindByID(c context.Context, id string) (domain.Author, error) {
	var err error
	var author domain.Author

	err = r.db.WithTransaction(c, func(tx pgx.Tx) error {
		author, err = r.AuthorQuery.FindByID(c, tx, id)
		return err
	})

	return author, err
}

func (r *authorRepository) ValidateAuthorName(c context.Context, name string) (domain.ValidateAuthorName, error) {
	var err error
	var author domain.ValidateAuthorName

	err = r.db.WithTransaction(c, func(tx pgx.Tx) error {
		author, err = r.AuthorQuery.ValidateAuthorName(c, tx, name)
		return err
	})

	return author, err
}

func (r *authorRepository) FindAllAuthor(c context.Context) ([]domain.Author, error) {
	var err error
	var authors []domain.Author

	err = r.db.WithTransaction(c, func(tx pgx.Tx) error {
		authors, err = r.AuthorQuery.FindAllAuthor(c, tx)
		return err
	})

	return authors, err
}

func (r *authorRepository) DeleteAuthor(c context.Context, id string) error {
	var err error

	err = r.db.WithTransaction(c, func(tx pgx.Tx) error {
		err = r.AuthorQuery.DeleteAuthor(c, tx, id)
		return err
	})

	return err
}
