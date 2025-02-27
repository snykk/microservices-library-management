package datatransfers

type BookRequest struct {
	Title      string `json:"title" validate:"required,min=3,max=255"`
	AuthorId   string `json:"author_id" validate:"required,uuid4"`
	CategoryId string `json:"category_id" validate:"required,uuid4"`
	Stock      int    `json:"stock" validate:"required,min=0"`
}

type BookUpdateRequest struct {
	Title      string `json:"title" validate:"required,min=3,max=255"`
	AuthorId   string `json:"author_id" validate:"required,uuid4"`
	CategoryId string `json:"category_id" validate:"required,uuid4"`
	Stock      int    `json:"stock" validate:"required,min=0"`
	Version    int    `json:"version" validate:"required,min=1"`
}

type BookDeleteRequest struct {
	Version int `json:"version" validate:"required:min=1"`
}
