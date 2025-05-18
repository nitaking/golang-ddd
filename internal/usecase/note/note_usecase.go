package usecase

import (
	"context"
	"go-clean-architecture-boilerplate/internal/domain/note"
)

// CreateNote DTO

type CreateNoteInput struct {
	Title, Content string
}
type CreateNoteOutput struct {
	ID note.NoteID
}

// EditNoteInput EditNote DTO
type EditNoteInput struct {
	ID      note.NoteID
	Title   string
	Content string
}

type EditNoteOutput struct {
	ID note.NoteID
}

// DeleteNoteInput DeleteNote DTO
type DeleteNoteInput struct {
	ID note.NoteID
}
type DeleteNoteOutput struct {
	ID note.NoteID
}

// LinkNoteInput LinkNote DTO
type LinkNoteInput struct {
	ID    note.NoteID
	Links []note.Link
}
type LinkNoteOutput struct {
	ID note.NoteID
}

type NoteUseCase interface {
	CreateNote(ctx context.Context, input CreateNoteInput) (*CreateNoteOutput, error)
	EditNote(ctx context.Context, input EditNoteInput) (*EditNoteOutput, error)
	DeleteNote(ctx context.Context, input DeleteNoteInput) (*DeleteNoteOutput, error)
	LinkNote(ctx context.Context, input LinkNoteInput) (*LinkNoteOutput, error)
	SearchNote(ctx context.Context, keyword string) (note.NoteSummary, error)
}

type noteUseCase struct {
	NoteRepository  note.NoteRepository
	QueryRepository note.NoteQueryRepository
}

func NewNoteUseCase(
	repo note.NoteRepository,
	queryRepo note.NoteQueryRepository,
) NoteUseCase {

	return &noteUseCase{
		NoteRepository:  repo,
		QueryRepository: queryRepo,
	}
}

func (n noteUseCase) CreateNote(ctx context.Context, input CreateNoteInput) (*CreateNoteOutput, error) {
	newNote, err := note.NewNote(input.Title, input.Content)

	if err != nil {
		return nil, err
	}
	err = n.NoteRepository.Save(ctx, newNote)
	if err != nil {
		return nil, err
	}
	return &CreateNoteOutput{
		ID: newNote.ID,
	}, nil
}

func (n noteUseCase) EditNote(ctx context.Context, in EditNoteInput) (*EditNoteOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (n noteUseCase) DeleteNote(ctx context.Context, in DeleteNoteInput) (*DeleteNoteOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (n noteUseCase) LinkNote(ctx context.Context, in LinkNoteInput) (*LinkNoteOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (n noteUseCase) SearchNote(ctx context.Context, keyword string) (note.NoteSummary, error) {
	result, err := n.QueryRepository.Search(ctx, keyword)
	if err != nil {
		return result, err
	}

	return result, nil
}
