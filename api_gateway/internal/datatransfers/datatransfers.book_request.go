package datatransfers

type BookRequest struct {
	Title      string `json:"title" validate:"required,min=3,max=255"`
	AuthorId   string `json:"author_id" validate:"required,uuid4"`
	CategoryId string `json:"category_id" validate:"required,uuid4"`
	Stock      int    `json:"stock" validate:"required,min=0"`
}
