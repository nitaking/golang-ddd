package bolt

import (
	"context"
	"encoding/json"
	"go-clean-architecture-boilerplate/internal/domain/note"
	"go.etcd.io/bbolt"
)

type NoteRepository struct {
	db *bbolt.DB
}

// NewNoteRepository creates a new NoteRepository instance.
func NewNoteRepository(db *bbolt.DB) *NoteRepository {
	return &NoteRepository{db: db}
}

func (r *NoteRepository) Save(ctx context.Context, n *note.Note) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(NoteBucketName))
		if err != nil {
			return err
		}
		data, err := json.Marshal(n)
		if err != nil {
			return err
		}
		return b.Put([]byte(n.ID), data)
	})
}

func (r *NoteRepository) FindByID(ctx context.Context, id note.NoteID) (*note.Note, error) {
	result := &note.Note{}
	err := r.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(NoteBucketName))
		if b == nil {
			return note.ErrNotFound
		}
		data := b.Get([]byte(id))
		if data == nil {
			return note.ErrNotFound
		}
		return json.Unmarshal(data, result)
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *NoteRepository) FindAll(ctx context.Context) ([]*note.Note, error) {
	//TODO implement me
	panic("implement me")
}

func (r *NoteRepository) DeleteByID(ctx context.Context, id note.NoteID) error {
	err := r.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(NoteBucketName))
		if b == nil {
			return note.ErrNotFound
		}
		return b.Delete([]byte(id))
	})
	if err != nil {
		return err
	}
	return nil
}
