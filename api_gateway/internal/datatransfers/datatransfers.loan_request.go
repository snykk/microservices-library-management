package datatransfers

type LoanRequest struct {
	BookId      string `json:"book_id" validate:"required,uuid4"`
	BookVersion int    `json:"book_version" validate:"required,min=1"`
}

type LoanStatusUpdateRequest struct {
	Status  string `json:"status" validate:"required,oneof=BORROWED RETURNED OVERDUE LOST"`
	Version int    `json:"version" validate:"required,min=1"`
}

type LoanReturnRequest struct {
	Version     int `json:"version" validate:"required,min=1"`
	BookVersion int `json:"book_version" validate:"required,min=1"`
}
