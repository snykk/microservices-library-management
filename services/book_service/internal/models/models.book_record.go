package models

import "time"

type BookRecord struct {
	Id         string    `db:"id"`
	Title      string    `db:"title"`
	AuthorId   string    `db:"author_id"`
	CategoryId string    `db:"category_id"`
	Stock      int       `db:"stock"`
	Version    int       `db:"version" json:"version"` // Field version untuk optimistic locking
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
