package domain

import (
	"test-backend-altech/model/web/response"

	"time"

	"github.com/google/uuid"
)

type Author struct {
	Id        string
	Name      string
	Bio       string
	BirthDate string
	CreatedAt time.Time
}

func (author *Author) GenerateID() {
	uuid := uuid.New().String()
	author.Id = uuid
}

func (j *Author) ToAuthorResponse() response.AuthorResponse {
	return response.AuthorResponse{
		Id:        j.Id,
		Name:      j.Name,
		Bio:       j.Bio,
		BirthDate: j.BirthDate,
	}
}

type ValidateAuthorName struct {
	Name string `json:"name"`
}

type UpdateAuthor struct {
	Name      string `json:"name"`
	Bio       string `json:"bio"`
	BirthDate string `json:"birth_date"`
}
