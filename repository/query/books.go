package query

import (
	"context"
	"log"
	"test-backend-altech/model/domain"
	"test-backend-altech/model/web/response"

	"github.com/jackc/pgx/v5"
)

type BookQuery interface {
	CreateBook(c context.Context, tx pgx.Tx, book domain.Book) error
	UpdateBook(c context.Context, tx pgx.Tx, id string, book domain.UpdateBook) error
	FindByID(c context.Context, tx pgx.Tx, id string) (response.BookResponse, error)
	ValidateBookTitle(c context.Context, tx pgx.Tx, id string) (domain.ValidateBookTitle, error)
	FindAllBook(c context.Context, tx pgx.Tx) ([]response.BookResponse, error)
	DeleteBook(c context.Context, tx pgx.Tx, id string) error
}

type BookQueryImpl struct {
}

func NewBook() BookQuery {
	return &BookQueryImpl{}
}

func (repository *BookQueryImpl) CreateBook(c context.Context, tx pgx.Tx, book domain.Book) error {

	query := `INSERT INTO books 
	(
		"id", 
		"title",
		"description",
		"author_id",
		"publish_date",
		"created_at"
	) 
	VALUES ($1,$2,$3,$4,$5,$6)`

	_, err := tx.Exec(c, query,
		book.Id,
		book.Title,
		book.Description,
		book.AuthorId,
		book.PublishDate,
		book.CreatedAt)

	return err
}

func (repository *BookQueryImpl) UpdateBook(c context.Context, tx pgx.Tx, id string, book domain.UpdateBook) error {
	query := `UPDATE  books SET 
	 title=$1,
	 description=$2,
	 publish_date=$3,
	 author_id=$4
	 WHERE id=$5
	`

	_, err := tx.Exec(c, query,
		book.Title,
		book.Description,
		book.PublishDate,
		book.AuthorId,
		id)

	return err
}
func (repository *BookQueryImpl) FindByID(c context.Context, tx pgx.Tx, id string) (response.BookResponse, error) {
	query := `
        SELECT
            b.id,
			b.title,
			b.description,
			b.publish_date,
			b.author_id,
			a.name
        FROM
            books AS b
		LEFT JOIN authors AS a ON b.author_id = a.id
        WHERE
            b.id = $1;
    `

	row := tx.QueryRow(c, query, id)
	var data response.BookResponse
	if err := row.Scan(
		&data.Id,
		&data.Title,
		&data.Description,
		&data.PublishDate,
		&data.AuthorId,
		&data.AuthorName,
	); err != nil {
		log.Println("Scan", err)
		return response.BookResponse{}, err
	}

	return data, nil
}

func (repository *BookQueryImpl) ValidateBookTitle(c context.Context, tx pgx.Tx, title string) (domain.ValidateBookTitle, error) {
	query := `
        SELECT
			b.title
        FROM
            books AS b
        WHERE
            b.title = $1;
    `

	row := tx.QueryRow(c, query, title)
	var data domain.ValidateBookTitle
	if err := row.Scan(
		&data.Title,
	); err != nil && err != pgx.ErrNoRows {
		log.Println("Scan", err)
		return domain.ValidateBookTitle{}, err
	}

	return data, nil
}

func (repository *BookQueryImpl) FindAllBook(c context.Context, tx pgx.Tx) ([]response.BookResponse, error) {

	query :=
		`
		SELECT
			b.id,
			b.title,
			b.description,
			b.publish_date,
			b.author_id,
			a.name
		FROM books AS b
		LEFT JOIN authors AS a ON b.author_id = a.id`

	rows, err := tx.Query(c, query)
	if err != nil {
		return nil, err
	}

	var datas []response.BookResponse
	for rows.Next() {
		var data response.BookResponse
		err := rows.Scan(&data.Id,
			&data.Title,
			&data.Description,
			&data.PublishDate,
			&data.AuthorId,
			&data.AuthorName)
		if err != nil {
			return nil,
				err
		}
		datas = append(datas, data)
	}

	defer rows.Close()
	return datas, nil
}

func (repository *BookQueryImpl) DeleteBook(c context.Context, tx pgx.Tx, id string) error {

	query := `DELETE FROM books WHERE id = $1`

	_, err := tx.Exec(c, query, id)
	if err != nil {
		return err
	}

	return nil
}
