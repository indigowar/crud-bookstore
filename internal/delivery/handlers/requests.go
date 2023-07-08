package handlers

import "net/http"

type BookDataRequest struct {
	Title           string `json:"title"`
	Author          string `json:"author"`
	PublicationYear int    `json:"publication_year"`
	ISBN            int    `json:"isbn"`
}

func (request *BookDataRequest) Bind(r *http.Request) error {
	return nil
}
