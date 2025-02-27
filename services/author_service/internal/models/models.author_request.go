package models

type AuthorCreateRequest struct {
	Name      string
	Biography string
}

type AuthorUpdateRequest struct {
	Name      string
	Biography string
	Version   int
}
