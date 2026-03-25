package core

import "errors"

var (
	ErrNotFound      = errors.New("resource not found")
	ErrDuplicateBook = errors.New("book already exists")
	ErrInvalidInput  = errors.New("invalid input data")
)
