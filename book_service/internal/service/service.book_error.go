package service

import "errors"

var (
	ErrGetBook     = errors.New("failed to retrieve book data")
	ErrGetListBook = errors.New("failed to retrieve book list")
	ErrUpdateBook  = errors.New("failed to update book data")
	ErrDeleteBook  = errors.New("failed to delete book data")
)
