package datatransfers

import "time"

type AuthorResponse struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Biography string    `json:"biography"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
