package bolt

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/ahmetb/go-linq/v3"
	"go-clean-architecture-boilerplate/internal/domain/note"
	"go.etcd.io/bbolt"
	"strings"
)

type NoteQueryRepository struct {
	db *bbolt.DB
}

func NewNoteQueryRepository(db *bbolt.DB) *NoteQueryRepository {
	return &NoteQueryRepository{db: db}
}

func (r *NoteQueryRepository) Search(ctx context.Context, query string) (note.NoteSummary, error) {
	// Performance is suboptimal here, as this is implemented using go-linq to retrieve all records and then filter with a Where clause
	// Consider implementing the query using SQL in a relational database or through a Search Engine approach

	var allNotes []note.Note

	// まずbboltから全件取得
	err := r.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("notes"))
		if b == nil {
			return errors.New("bucket not found")
		}
		return b.ForEach(func(k, v []byte) error {
			var n note.Note
			if err := json.Unmarshal(v, &n); err != nil {
				return err
			}
			allNotes = append(allNotes, n)
			return nil
		})
	})
	if err != nil {
		return note.NoteSummary{}, err
	}

	// go-linqでWhereフィルタ
	var result []note.Note
	linq.From(allNotes).
		Where(func(item interface{}) bool {
			n := item.(note.Note)
			return strings.Contains(n.Title, query) || strings.Contains(n.Content, query)
		}).
		Select(func(item interface{}) interface{} {
			n := item.(note.Note)
			return n
		}).
		ToSlice(&result)

	return note.NoteSummary{
		Count: len(result),
		Notes: result,
	}, nil
}
