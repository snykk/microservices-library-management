package datatransfers

import "time"

type CategoryResponse struct {
	Id          string          `json:"id"`
	Name        string          `json:"name"`
	SampleBooks *[]BookResponse `json:"sample_books,omitempty"` // Sample of book (previous scenario no longer valid because now book service use pagination)
	TotalBooks  *int            `json:"total_books,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
