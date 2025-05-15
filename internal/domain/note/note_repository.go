package note

import (
	"context"
)

type NoteRepository interface {
	Save(ctx context.Context, h *Note) error
	DeleteByID(ctx context.Context, id NoteID) error
	FindByID(ctx context.Context, id NoteID) (*Note, error)
	FindAll(ctx context.Context) ([]*Note, error)
}
