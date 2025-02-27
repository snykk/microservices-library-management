package datatransfers

type CategoryRequest struct {
	Name string `json:"name" validate:"required,min=3,max=100"`
}

type CategoryUpdateRequest struct {
	Name    string `json:"name" validate:"required,min=3,max=100"`
	Version int    `json:"version" validate:"required,min=1"`
}
