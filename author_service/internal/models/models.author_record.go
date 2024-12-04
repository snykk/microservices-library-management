package models

import "time"

type AuthorRecord struct {
	Id        string    `db:"id"`
	Name      string    `db:"name"`
	Biography string    `db:"biography"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
