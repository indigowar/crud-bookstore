package models

// Book - title, author, publication year, and ISBN.
type Book struct {
	ID              int64  `json:"id"`
	Title           string `json:"title"`
	Author          string `json:"author"`
	PublicationYear int    `json:"publication_year"`
	ISBN            int    `json:"isbn"`
}
