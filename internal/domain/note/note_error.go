package note

import "errors"

var (
	ErrNotFound   = errors.New("note: not found")
	ErrEmptyTitle = errors.New("note: title is empty")
)
