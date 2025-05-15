package note

import "context"

type NoteService interface {
	SaveNote(ctx context.Context, h *Note) error
}
