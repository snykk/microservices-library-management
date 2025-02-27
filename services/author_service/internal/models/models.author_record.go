package models

import "time"

type AuthorRecord struct {
	Id        string    `db:"id" json:"id"` // json tags for logging (publish rabbitmq)
	Name      string    `db:"name" json:"name"`
	Biography string    `db:"biography" json:"biography"`
	Version   int       `db:"version" json:"version"` // Field version untuk optimistic locking
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
