package models

// Book - title, author, publication year, and ISBN.
type Book struct {
	ID              int64
	Title           string
	Author          string
	PublicationYear int
	ISBN            int
}
