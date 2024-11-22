package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"test-backend-altech/config"
	"test-backend-altech/exception"
	"test-backend-altech/model/domain"
	request "test-backend-altech/model/web/req"
	response "test-backend-altech/model/web/response"
	"test-backend-altech/repository"
)

type BookService interface {
	CreateBook(ctx context.Context, request request.BookRequest) (response.BookResponse, error)
	FindByID(ctx context.Context, id string) (response.BookResponse, error)
	UpdateBook(ctx context.Context, request request.BookRequest, id string) (response.BookResponse, error)
	FindAllBook(ctx context.Context, cache config.Cache) ([]response.BookResponse, error)
	DeleteBook(ctx context.Context, id string) (response.BookResponse, error)
}

type bookService struct {
	bookRepository repository.BookRepository
	cache          config.Cache
}

func NewBookService(bookRepository repository.BookRepository, cache config.Cache) BookService {
	return &bookService{
		bookRepository: bookRepository,
		cache:          cache,
	}
}

func (s *bookService) CreateBook(c context.Context, request request.BookRequest) (response.BookResponse, error) {
	book := domain.Book{
		Title:       request.Title,
		Description: request.Description,
		PublishDate: request.PublishDate,
		AuthorId:    request.AuthorId,
		CreatedAt:   time.Now(),
	}
	book.GenerateID()

	validateTitle, err := s.bookRepository.ValidateBookTitle(c, book.Title)
	if err != nil {
		return response.BookResponse{}, err
	}
	if validateTitle.Title != "" {
		return response.BookResponse{}, exception.ErrValidateBadRequest("Book name already exists", request)
	}

	if err := s.bookRepository.CreateBook(c, book); err != nil {
		return response.BookResponse{}, err
	}

	newBook, err := s.bookRepository.FindByID(c, book.Id)
	if err != nil {
		return response.BookResponse{}, exception.ErrInternalServer(fmt.Sprintf("Successfully created book, but failed to get the created book. Error: %s", err.Error()))
	}

	return newBook, err
}
func (s *bookService) FindByID(ctx context.Context, id string) (response.BookResponse, error) {
	res, err := s.bookRepository.FindByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return response.BookResponse{}, exception.ErrNotFound("Book not found")
		} else {
			return response.BookResponse{}, err
		}
	}

	return res, err
}

func (s *bookService) UpdateBook(ctx context.Context, request request.BookRequest, id string) (response.BookResponse, error) {
	data, err := s.bookRepository.FindByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return response.BookResponse{}, exception.ErrNotFound("Book not found")
		} else {
			return response.BookResponse{}, err
		}
	}

	book := domain.UpdateBook{
		Title:       request.Title,
		Description: request.Description,
		PublishDate: request.PublishDate,
		AuthorId:    request.AuthorId,
	}
	data.Description = book.Description
	data.Title = book.Title
	data.PublishDate = book.PublishDate
	data.AuthorId = book.AuthorId

	if err := s.bookRepository.UpdateBook(ctx, id, book); err != nil {
		return response.BookResponse{}, err
	}
	return data, err
}

func (s *bookService) FindAllBook(ctx context.Context, cache config.Cache) ([]response.BookResponse, error) {
	bookKey := "books"
	cacheKey := fmt.Sprintf("driver:%s", bookKey)

	var dataRedis []response.BookResponse
	cachedData, err := s.cache.Get(ctx, cacheKey)
	if err == nil {
		if err := json.Unmarshal(cachedData, &dataRedis); err == nil {
			log.Println("Cache hit")
			return dataRedis, nil
		} else {
			log.Printf("Failed to unmarshal cached data: %v", err)
		}
	}

	res, err := s.bookRepository.FindAllBook(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return []response.BookResponse{}, exception.ErrNotFound("Book not found")
		} else {
			return []response.BookResponse{}, err
		}
	}

	var data []response.BookResponse

	data = append(data, res...)

	dbResponseBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	if err := s.cache.Set(ctx, cacheKey, dbResponseBytes, time.Hour*168); err != nil {
		log.Printf("Failed to cache data: %v", err)
	}
	return data, err
}

func (s *bookService) DeleteBook(ctx context.Context, id string) (response.BookResponse, error) {
	data, err := s.bookRepository.FindByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return response.BookResponse{}, exception.ErrNotFound("Book not found")
		} else {
			return response.BookResponse{}, err
		}
	}

	if err := s.bookRepository.DeleteBook(ctx, id); err != nil {
		return response.BookResponse{}, err
	}
	return data, err
}
