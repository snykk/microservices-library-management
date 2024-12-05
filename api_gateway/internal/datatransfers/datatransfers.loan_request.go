package datatransfers

import "time"

type LoanRequest struct {
	BookId string `json:"book_id"`
}

type LoanStatusUpdateRequest struct {
	Status     string    `json:"status"`
	ReturnDate time.Time // assign in handler
}
