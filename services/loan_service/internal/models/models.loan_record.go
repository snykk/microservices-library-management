package models

import "time"

type LoanRecord struct {
	Id         string     `db:"id"`
	UserId     string     `db:"user_id"`
	BookId     string     `db:"book_id"`
	LoanDate   time.Time  `db:"loan_date"`
	ReturnDate *time.Time `db:"return_date"`
	Status     string     `db:"status"` // "BORROWED", "RETURNED", "LOST"
	Version    int        `db:"version"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
}
