package memory

import (
	"context"
	"go-clean-architecture-boilerplate/internal/domain/note"
	"strings"
	"sync"
)

type NoteQueryRepository struct {
	mu   sync.RWMutex
	data map[note.NoteID]*note.Note
}

func NewNoteQueryRepository(data map[note.NoteID]*note.Note) *NoteQueryRepository {
	return &NoteQueryRepository{
		data: data,
	}
}

func (r *NoteQueryRepository) Search(ctx context.Context, query string) ([]note.NoteSummary, error) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []note.NoteSummary

	for _, n := range r.data {
		if strings.Contains(strings.ToLower(n.Title), strings.ToLower(query)) {
			result = append(result, note.NoteSummary{
				ID:    n.ID,
				Title: n.Title,
			})
		}
	}
	return result, nil
}
