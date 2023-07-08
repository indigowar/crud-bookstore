package services

import (
	"bookstore/internal/domain/models"
	"context"
	"errors"
	"fmt"
	"golang.org/x/exp/slog"
	"math/rand"
	"strconv"
	"time"
)

var (
	ErrBookDoesNotExist  = errors.New("book does not exist")
	ErrBookAlreadyExists = errors.New("book is already exists")
	ErrInternalError     = errors.New("internal error")
	ErrInvalidData       = errors.New("invalid data")
)

type BookStorage interface {
	GetByID(ctx context.Context, id int64) (models.Book, error)
	GetByISBN(ctx context.Context, isbn int) (models.Book, error)
	GetByTitleAndAuthor(ctx context.Context, author string, title string) (models.Book, error)

	GetAll(ctx context.Context) ([]models.Book, error)

	Add(ctx context.Context, book models.Book) error
	Remove(ctx context.Context, id int64) error

	Update(ctx context.Context, book models.Book) error
}

type BookService interface {
	GetAll(ctx context.Context) ([]models.Book, error)
	GetByID(ctx context.Context, bookID int64) (models.Book, error)
	Create(ctx context.Context, title string, author string, publicationYear int, isbn int) (models.Book, error)

	Update(ctx context.Context, book models.Book) error
	Delete(ctx context.Context, bookID int64) error
}

type bookService struct {
	storage BookStorage
	logger  *slog.Logger
}

func NewBookService(logger *slog.Logger, storage BookStorage) BookService {
	return &bookService{
		storage: storage,
		logger:  logger,
	}
}

func (svc bookService) GetAll(ctx context.Context) ([]models.Book, error) {
	result, err := svc.storage.GetAll(ctx)
	if err != nil {
		svc.logger.Warn("BookService.GetAll() -> %w", err)
		return nil, ErrInternalError
	}
	return result, nil
}

func (svc bookService) GetByID(ctx context.Context, bookID int64) (models.Book, error) {
	result, err := svc.storage.GetByID(ctx, bookID)
	if err == nil {
		return result, nil
	}

	if errors.Is(err, ErrBookDoesNotExist) {
		return models.Book{}, fmt.Errorf("%d %w", bookID, ErrBookDoesNotExist)
	}
	return models.Book{}, ErrInternalError
}

func (svc bookService) Create(ctx context.Context, title string, author string, publicationYear int, isbn int) (models.Book, error) {
	if len(title) < 2 {
		return models.Book{}, fmt.Errorf("title: %w", ErrInvalidData)
	}

	if len(author) < 2 {
		return models.Book{}, fmt.Errorf("author: %w", ErrInvalidData)
	}

	if publicationYear < 1500 || publicationYear > time.Now().Year() {
		return models.Book{}, fmt.Errorf("publication year: %w", ErrInvalidData)
	}

	if isbnInStr := strconv.Itoa(isbn); len(isbnInStr) != 13 {
		return models.Book{}, fmt.Errorf("ISBN: %w", ErrInvalidData)
	}

	if _, err := svc.storage.GetByTitleAndAuthor(ctx, title, author); err == nil {
		return models.Book{}, fmt.Errorf("%s by %s - %w", title, author, ErrBookAlreadyExists)
	}

	if _, err := svc.storage.GetByISBN(ctx, isbn); err == nil {
		return models.Book{}, fmt.Errorf("ISBN: %d - %w", isbn, ErrBookAlreadyExists)
	}

	book := models.Book{
		ID:              rand.Int63(),
		Title:           title,
		Author:          author,
		PublicationYear: publicationYear,
		ISBN:            isbn,
	}

	err := svc.storage.Add(ctx, book)
	if err != nil {
		svc.logger.Warn("BookService.Create -> %w", err)
		return models.Book{}, ErrInternalError
	}

	return book, nil
}

func (svc bookService) Update(ctx context.Context, book models.Book) error {
	if len(book.Title) < 2 {
		return fmt.Errorf("title: %w", ErrInvalidData)
	}

	if len(book.Author) < 2 {
		return fmt.Errorf("author: %w", ErrInvalidData)
	}

	if book.PublicationYear < 1500 || book.PublicationYear > time.Now().Year() {
		return fmt.Errorf("publication year: %w", ErrInvalidData)
	}

	if isbnInStr := strconv.Itoa(book.ISBN); len(isbnInStr) != 13 {
		return fmt.Errorf("ISBN: %w", ErrInvalidData)
	}

	if err := svc.storage.Update(ctx, book); err != nil {
		if errors.Is(err, ErrBookDoesNotExist) {
			return fmt.Errorf("%d %w", book.ID, ErrBookDoesNotExist)
		}
		svc.logger.Warn("BookService.Update -> %w", err)
		return ErrInternalError
	}
	return nil
}

func (svc bookService) Delete(ctx context.Context, bookID int64) error {
	err := svc.storage.Remove(ctx, bookID)
	if err == nil {
		return nil
	}

	if errors.Is(err, ErrBookDoesNotExist) {
		return fmt.Errorf("%d %w", bookID, ErrBookDoesNotExist)
	}
	svc.logger.Warn("BookService.Delete -> %w", err)
	return ErrInternalError
}
