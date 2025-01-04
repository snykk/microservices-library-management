package service

import "errors"

var (
	ErrCreateCategory  = errors.New("failed to create new category")
	ErrGetCategory     = errors.New("failed to retrieve category data")
	ErrGetListCategory = errors.New("failed to retrieve category list")
	ErrUpdateCategory  = errors.New("failed to update category data")
	ErrDeleteCategory  = errors.New("failed to delete category data")
)
