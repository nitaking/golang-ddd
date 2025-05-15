package memory

import (
	"context"
	"go-clean-architecture-boilerplate/internal/domain/note"
	"sync"
	"time"
)

// NoteRepository は note.Repository を実装するメモリリポジトリ
type NoteRepository struct {
	mu   sync.RWMutex
	data map[note.NoteID]*note.Note
}

func NewNoteRepository(data map[note.NoteID]*note.Note) *NoteRepository {
	return &NoteRepository{
		data: data,
	}
}

// Save は Note の新規作成 or 更新
func (r *NoteRepository) Save(ctx context.Context, n *note.Note) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	// タイムスタンプ未設定なら設定
	if n.CreatedAt.IsZero() {
		n.CreatedAt = time.Now()
	}
	n.UpdatedAt = time.Now()
	r.data[n.ID] = n
	return nil
}

// FindByID は ID での取得
func (r *NoteRepository) FindByID(ctx context.Context, id note.NoteID) (*note.Note, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	n, ok := r.data[id]
	if !ok {
		return nil, note.ErrNotFound
	}
	return n, nil
}

func (r *NoteRepository) DeleteByID(ctx context.Context, id note.NoteID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.data[id]
	if !ok {
		return note.ErrNotFound
	}
	delete(r.data, id)
	return nil
}

func (r *NoteRepository) FindAll(ctx context.Context) ([]*note.Note, error) {
	//TODO implement me
	panic("implement me")
}
