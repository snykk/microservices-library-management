package datatransfers

type AuthorRequest struct {
	Name      string `json:"name" validate:"required,min=3,max=100"`
	Biography string `json:"biography" validate:"required"`
	Version   int    `json:"version" validate:"required"`
}

type AuthorDeleteRequest struct {
	Version int `json:"version" validate:"required"`
}
