package datatransfers

type AuthorRequest struct {
	Name      string `json:"name" validate:"required,min=3,max=100"`
	Biography string `json:"biography" validate:"required"`
}
