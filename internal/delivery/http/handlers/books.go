package handlers

import (
	"bookstore/internal/services"
	"net/http"
)

func MakeGetSpecificBookHandler(svc services.BookService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func MakeGetAllBooksHandler(svc services.BookService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func MakeCreateBookHandler(svc services.BookService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func MakeDeleteBookHandler(svc services.BookService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func MakeUpdateBookHandler(svc services.BookService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

type creationInfo struct {
	Title           string `json:"title"`
	Author          string `json:"author"`
	PublicationYear int    `json:"publication_year"`
	ISBN            int    `json:"isbn"`
}
