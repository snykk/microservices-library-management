package datatransfers

type CategoryRequest struct {
	Name string `json:"name" validate:"required,min=3,max=100"`
}
