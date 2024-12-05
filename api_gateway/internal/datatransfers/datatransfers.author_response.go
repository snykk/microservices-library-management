package datatransfers

import "time"

type AuthorResponse struct {
	Id        string          `json:"id"`
	Name      string          `json:"name"`
	Biography string          `json:"biography"`
	Books     *[]BookResponse `json:"books,omitempty"`
	CreatedAt time.Time       `json:"created_at,omitempty"`
	UpdatedAt time.Time       `json:"updated_at,omitempty"`
}
