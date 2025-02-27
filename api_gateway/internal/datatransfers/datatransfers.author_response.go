package datatransfers

import "time"

type AuthorResponse struct {
	Id          string          `json:"id"`
	Name        string          `json:"name"`
	Biography   string          `json:"biography"`
	SampleBooks *[]BookResponse `json:"sample_books,omitempty"` // Sample of book (previous scenario no longer valid because now book service use pagination)
	TotalBooks  *int            `json:"total_books,omitempty"`
	Version     int             `json:"version"`
	CreatedAt   time.Time       `json:"created_at,omitempty"`
	UpdatedAt   time.Time       `json:"updated_at,omitempty"`
}
