package handlers

import (
	"bookstore/internal/domain/models"
	"github.com/go-chi/render"
	"net/http"
)

type BookResponse struct {
	models.Book

	Elapsed int64 `json:"elapsed"`
}

func (br *BookResponse) Render(w http.ResponseWriter, r *http.Request) error {
	br.Elapsed = 10
	return nil
}

func NewBookResponse(book models.Book) *BookResponse {
	return &BookResponse{Book: book}
}

func NewBookListResponse(books []models.Book) []render.Renderer {
	result := make([]render.Renderer, len(books))
	for i, v := range books {
		result[i] = NewBookResponse(v)
	}
	return result
}

type ErrorResponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	StatusText string `json:"status"`
	AppCode    int64  `json:"code,omitempty"`
	ErrorText  string `json:"error,omitempty"`
}

func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrorResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrorResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnprocessableEntity,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

func ErrInternalError(err error) render.Renderer {
	return &ErrorResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "Internal Server Error.",
		ErrorText:      err.Error(),
	}
}

var ErrNotFound = &ErrorResponse{HTTPStatusCode: http.StatusNotFound, StatusText: "Resource not found."}
