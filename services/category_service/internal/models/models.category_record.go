package models

import "time"

type CategoryRecord struct {
	Id        string    `db:"id"`
	Name      string    `db:"name"`
	Version   int       `db:"version"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
