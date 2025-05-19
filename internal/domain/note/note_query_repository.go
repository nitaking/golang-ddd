package note

import (
	"context"
)

type NoteQueryRepository interface {
	Search(ctx context.Context, query string) (NoteSummary, error)
}
