package datatransfers

import "time"

type LoanResponse struct {
	Id         string     `json:"id"`
	UserId     string     `json:"user_id"`
	BookId     string     `json:"book_id"`
	LoanDate   time.Time  `json:"loan_date"`
	ReturnDate *time.Time `json:"return_date"`
	Status     string     `json:"status"`
	Version    int        `json:"version"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
