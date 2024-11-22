package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"test-backend-altech/exception"
	"test-backend-altech/model/domain"
	request "test-backend-altech/model/web/req"
	response "test-backend-altech/model/web/response"
	"test-backend-altech/repository"
)

type AuthorService interface {
	CreateAuthor(ctx context.Context, request request.AuthorRequest) (response.AuthorResponse, error)
	FindByID(ctx context.Context, id string) (response.AuthorResponse, error)
	UpdateAuthor(ctx context.Context, request request.AuthorRequest, id string) (response.AuthorResponse, error)
	FindAllAuthor(ctx context.Context) ([]response.AuthorResponse, error)
	DeleteAuthor(ctx context.Context, id string) (response.AuthorResponse, error)
}

type authorService struct {
	authorRepository repository.AuthorRepository
}

func NewAuthorService(authorRepository repository.AuthorRepository) AuthorService {
	return &authorService{
		authorRepository: authorRepository,
	}
}

func (s *authorService) CreateAuthor(c context.Context, request request.AuthorRequest) (response.AuthorResponse, error) {
	author := domain.Author{
		Name:      request.Name,
		Bio:       request.Bio,
		BirthDate: request.BirthDate,
		CreatedAt: time.Now(),
	}
	author.GenerateID()

	validateName, err := s.authorRepository.ValidateAuthorName(c, author.Name)
	if err != nil {
		return response.AuthorResponse{}, err
	}
	if validateName.Name != "" {
		return response.AuthorResponse{}, exception.ErrValidateBadRequest("Author name already exists", request)
	}

	if err := s.authorRepository.CreateAuthor(c, author); err != nil {
		return response.AuthorResponse{}, err
	}

	newAuthor, err := s.authorRepository.FindByID(c, author.Id)
	if err != nil {
		return response.AuthorResponse{}, exception.ErrInternalServer(fmt.Sprintf("Successfully created author, but failed to get the created author. Error: %s", err.Error()))
	}

	return newAuthor.ToAuthorResponse(), err
}
func (s *authorService) FindByID(ctx context.Context, id string) (response.AuthorResponse, error) {
	res, err := s.authorRepository.FindByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return response.AuthorResponse{}, exception.ErrNotFound("Author not found")
		} else {
			return response.AuthorResponse{}, err
		}
	}

	return res.ToAuthorResponse(), err
}

func (s *authorService) UpdateAuthor(ctx context.Context, request request.AuthorRequest, id string) (response.AuthorResponse, error) {
	data, err := s.authorRepository.FindByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return response.AuthorResponse{}, exception.ErrNotFound("Author not found")
		} else {
			return response.AuthorResponse{}, err
		}
	}

	author := domain.UpdateAuthor{
		Name:      request.Name,
		Bio:       request.Bio,
		BirthDate: request.BirthDate,
	}
	data.Bio = author.Bio
	data.Name = author.Name
	data.BirthDate = author.BirthDate

	if err := s.authorRepository.UpdateAuthor(ctx, id, author); err != nil {
		return response.AuthorResponse{}, err
	}
	return data.ToAuthorResponse(), err
}

func (s *authorService) FindAllAuthor(ctx context.Context) ([]response.AuthorResponse, error) {
	res, err := s.authorRepository.FindAllAuthor(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return []response.AuthorResponse{}, exception.ErrNotFound("Author not found")
		} else {
			return []response.AuthorResponse{}, err
		}
	}

	var data []response.AuthorResponse
	for _, v := range res {
		data = append(data, v.ToAuthorResponse())
	}
	return data, err
}

func (s *authorService) DeleteAuthor(ctx context.Context, id string) (response.AuthorResponse, error) {
	data, err := s.authorRepository.FindByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return response.AuthorResponse{}, exception.ErrNotFound("Author not found")
		} else {
			return response.AuthorResponse{}, err
		}
	}

	if err := s.authorRepository.DeleteAuthor(ctx, id); err != nil {
		return response.AuthorResponse{}, err
	}
	return data.ToAuthorResponse(), err
}
