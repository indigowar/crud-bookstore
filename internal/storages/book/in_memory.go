package book

import (
	"bookstore/internal/domain/models"
	"bookstore/internal/domain/services"
	"context"
)

type storage struct {
	books map[int64]models.Book
}

func NewInMemoryBookStorage() services.BookStorage {
	return &storage{
		books: make(map[int64]models.Book),
	}
}

func (s storage) GetByID(_ context.Context, id int64) (models.Book, error) {
	b, ok := s.books[id]
	if !ok {
		return models.Book{}, services.ErrBookDoesNotExist
	}
	return b, nil
}

func (s storage) GetByISBN(_ context.Context, isbn int) (models.Book, error) {
	for _, v := range s.books {
		if v.ISBN == isbn {
			return v, nil
		}
	}
	return models.Book{}, services.ErrBookDoesNotExist
}

func (s storage) GetByTitleAndAuthor(_ context.Context, author string, title string) (models.Book, error) {
	for _, v := range s.books {
		if v.Author == author && v.Title == title {
			return v, nil
		}
	}
	return models.Book{}, services.ErrBookDoesNotExist
}

func (s storage) GetAll(_ context.Context) ([]models.Book, error) {
	result := make([]models.Book, len(s.books))
	i := 0
	for _, v := range s.books {
		result[i] = v
		i++
	}
	return result, nil
}

func (s storage) Add(_ context.Context, book models.Book) error {
	s.books[book.ID] = book
	return nil
}

func (s storage) Remove(_ context.Context, id int64) error {
	delete(s.books, id)
	return nil
}

func (s storage) Update(_ context.Context, book models.Book) error {
	s.books[book.ID] = book
	return nil
}
