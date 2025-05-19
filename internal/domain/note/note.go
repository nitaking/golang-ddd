package note

import (
	"errors"
	"github.com/google/uuid"
	"strings"
	"time"
)

// NoteID represents a unique identifier for a habit.
type NoteID string

func NewNoteID() NoteID {
	return NoteID(uuid.NewString())
}

type Note struct {
	ID      NoteID
	Title   string
	Content string
	Links   []Link

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewNote(title, content string) (*Note, error) {
	if strings.TrimSpace(title) == "" {
		return nil, errors.New("title cannot be empty")
	}

	return &Note{
		ID:        NewNoteID(),
		Title:     title,
		Content:   content,
		Links:     []Link{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

type ReconstructNoteInput struct {
	ID        NoteID
	Title     string
	Content   string
	Links     []Link
	CreatedAt time.Time
	UpdatedAt time.Time
}

func Reconstruct(input ReconstructNoteInput) (*Note, error) {
	if strings.TrimSpace(input.Title) == "" {
		return nil, errors.New("title cannot be empty")
	}

	return &Note{
		ID:        input.ID,
		Title:     input.Title,
		Content:   input.Content,
		Links:     input.Links,
		CreatedAt: input.CreatedAt,
		UpdatedAt: input.UpdatedAt,
	}, nil
}
