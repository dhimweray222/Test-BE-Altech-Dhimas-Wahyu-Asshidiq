package query

import (
	"context"
	"log"
	"test-backend-altech/model/domain"

	"github.com/jackc/pgx/v5"
)

type AuthorQuery interface {
	CreateAuthor(c context.Context, tx pgx.Tx, author domain.Author) error
	UpdateAuthor(c context.Context, tx pgx.Tx, id string, author domain.UpdateAuthor) error
	FindByID(c context.Context, tx pgx.Tx, id string) (domain.Author, error)
	ValidateAuthorName(c context.Context, tx pgx.Tx, id string) (domain.ValidateAuthorName, error)
	FindAllAuthor(c context.Context, tx pgx.Tx) ([]domain.Author, error)
	DeleteAuthor(c context.Context, tx pgx.Tx, id string) error
}

type AuthorQueryImpl struct {
}

func NewAuthor() AuthorQuery {
	return &AuthorQueryImpl{}
}

func (repository *AuthorQueryImpl) CreateAuthor(c context.Context, tx pgx.Tx, author domain.Author) error {

	query := `INSERT INTO authors 
	(
		"id", 
		"name",
		"bio",
		"birth_date",
		"created_at"
	) 
	VALUES ($1,$2,$3,$4,$5)`

	_, err := tx.Exec(c, query,
		author.Id,
		author.Name,
		author.Bio,
		author.BirthDate,
		author.CreatedAt)

	return err
}

func (repository *AuthorQueryImpl) UpdateAuthor(c context.Context, tx pgx.Tx, id string, author domain.UpdateAuthor) error {
	query := `UPDATE  authors SET 
	 name=$1,
	 bio=$2,
	 birth_date=$3
	 WHERE id=$4
	`

	_, err := tx.Exec(c, query,
		author.Name,
		author.Bio,
		author.BirthDate,
		id)

	return err
}
func (repository *AuthorQueryImpl) FindByID(c context.Context, tx pgx.Tx, id string) (domain.Author, error) {
	query := `
        SELECT
            a.id,
			a.name,
			a.bio,
			a.birth_date,
			a.created_at
        FROM
            authors AS a
        WHERE
            a.id = $1;
    `

	row := tx.QueryRow(c, query, id)
	var data domain.Author
	if err := row.Scan(
		&data.Id,
		&data.Name,
		&data.Bio,
		&data.BirthDate,
		&data.CreatedAt,
	); err != nil {
		log.Println("Scan", err)
		return domain.Author{}, err
	}

	return data, nil
}

func (repository *AuthorQueryImpl) ValidateAuthorName(c context.Context, tx pgx.Tx, name string) (domain.ValidateAuthorName, error) {
	query := `
        SELECT
			a.name
        FROM
            authors AS a
        WHERE
            a.name = $1;
    `

	row := tx.QueryRow(c, query, name)
	var data domain.ValidateAuthorName
	if err := row.Scan(
		&data.Name,
	); err != nil && err != pgx.ErrNoRows {
		log.Println("Scan", err)
		return domain.ValidateAuthorName{}, err
	}

	return data, nil
}

func (repository *AuthorQueryImpl) FindAllAuthor(c context.Context, tx pgx.Tx) ([]domain.Author, error) {

	query :=
		`SELECT
		 a.id,
		 a.name,
		 a.bio,
		 a.birth_date
			FROM authors AS a`

	rows, err := tx.Query(c, query)
	if err != nil {
		return nil, err
	}

	var datas []domain.Author
	for rows.Next() {
		var data domain.Author
		err := rows.Scan(&data.Id, &data.Name, &data.Bio, &data.BirthDate)
		if err != nil {
			return nil, err
		}
		datas = append(datas, data)
	}

	defer rows.Close()
	return datas, nil
}

func (repository *AuthorQueryImpl) DeleteAuthor(c context.Context, tx pgx.Tx, id string) error {

	query := `DELETE FROM authors WHERE id = $1`

	_, err := tx.Exec(c, query, id)
	if err != nil {
		return err
	}

	return nil
}
