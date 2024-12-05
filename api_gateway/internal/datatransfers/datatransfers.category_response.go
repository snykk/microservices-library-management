package datatransfers

import "time"

type CategoryResponse struct {
	Id        string          `json:"id"`
	Name      string          `json:"name"`
	Books     *[]BookResponse `json:"books,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
