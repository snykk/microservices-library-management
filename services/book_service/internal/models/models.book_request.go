package models

type BookRequest struct {
	Title      string
	AuthorId   string
	CategoryId string
	Stock      int
	Version    int
}
