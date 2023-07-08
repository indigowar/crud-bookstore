package handlers

import (
	"bookstore/internal/domain/models"
	"bookstore/internal/domain/services"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
)

type BookEndpoint struct {
	svc services.BookService
}

func NewBookEndpoint(svc services.BookService) *BookEndpoint {
	return &BookEndpoint{
		svc: svc,
	}
}

func (e *BookEndpoint) GetSpecificBook(w http.ResponseWriter, r *http.Request) {
	bookId, err := e.getIDFromRequest(r)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	result, err := e.svc.GetByID(r.Context(), bookId)
	if err != nil {
		if errors.Is(err, services.ErrBookDoesNotExist) {
			_ = render.Render(w, r, ErrNotFound)
			return
		}
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	_ = render.Render(w, r, NewBookResponse(result))
	return
}

func (e *BookEndpoint) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := e.svc.GetAll(r.Context())
	if err != nil {
		_ = render.Render(w, r, ErrInternalError(err))
		return
	}

	_ = render.RenderList(w, r, NewBookListResponse(books))
	return
}

func (e *BookEndpoint) CreateBook(w http.ResponseWriter, r *http.Request) {
	data, err := e.parseCreationInfo(r)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	book, err := e.svc.Create(r.Context(), data.Title, data.Author, data.PublicationYear, data.ISBN)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusCreated)
	_ = render.Render(w, r, NewBookResponse(book))
}

func (e *BookEndpoint) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := e.getIDFromRequest(r)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err = e.svc.Delete(r.Context(), id); err != nil {
		return
	}

	render.Status(r, http.StatusOK)
}

func (e *BookEndpoint) UpdateBook(w http.ResponseWriter, r *http.Request) {
	id, err := e.getIDFromRequest(r)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	data, err := e.parseCreationInfo(r)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	book := models.Book{
		ID:              id,
		Title:           data.Title,
		Author:          data.Author,
		PublicationYear: data.PublicationYear,
		ISBN:            data.ISBN,
	}

	if err := e.svc.Update(r.Context(), book); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusOK)
}

// getIDFromRequest - gets the id from http.Request
func (e *BookEndpoint) getIDFromRequest(r *http.Request) (int64, error) {
	paramBookId := chi.URLParam(r, "id")
	if paramBookId == "" {
		return 0, errors.New("no id specified")
	}

	bookId, err := strconv.ParseInt(paramBookId, 10, 64)
	if err != nil {
		return 0, err
	}
	return bookId, nil
}

// parseCreationInfo - decode the request's body to BookDataRequest
func (e *BookEndpoint) parseCreationInfo(r *http.Request) (BookDataRequest, error) {
	data := BookDataRequest{}
	if err := render.Bind(r, &data); err != nil {
		return BookDataRequest{}, err
	}
	return data, nil
}
