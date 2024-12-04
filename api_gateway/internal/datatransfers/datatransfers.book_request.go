package datatransfers

type BookRequest struct {
	Title      string `json:"title"`
	AuthorId   string `json:"author_id"`
	CategoryId string `json:"category_id"`
	Stock      int    `json:"stock"`
}
