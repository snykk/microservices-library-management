package datatransfers

import "time"

type BookResponse struct {
	Id         string            `json:"id"`
	Title      string            `json:"title"`
	AuthorId   *string           `json:"author_id,omitempty"`
	Author     *AuthorResponse   `json:"author,omitempty"`
	CategoryId *string           `json:"category_id,omitempty"`
	Category   *CategoryResponse `json:"category,omitempty"`
	Stock      int               `json:"stock"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}
