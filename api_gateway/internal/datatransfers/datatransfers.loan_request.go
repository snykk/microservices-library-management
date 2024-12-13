package datatransfers

import "time"

type LoanRequest struct {
	BookId string `json:"book_id" validate:"required,uuid4"`
}

type LoanStatusUpdateRequest struct {
	Status     string `json:"status" validate:"required,oneof=BORROWED RETURNED OVERDUE LOST"`
	ReturnDate time.Time
}
