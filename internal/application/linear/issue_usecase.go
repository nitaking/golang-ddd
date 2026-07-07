package linear

import (
	"context"
	"errors"

	domain "go-clean-architecture-boilerplate/internal/domain/linear"
)

type IssueUseCase interface {
	CreateIssue(ctx context.Context, input CreateIssueInput) (*IssueOutput, error)
	CompleteIssue(ctx context.Context, input CompleteIssueInput) (*IssueOutput, error)
	ListIssues(ctx context.Context, input ListIssuesInput) ([]IssueOutput, error)
}

type CreateIssueInput struct {
	TeamID      string
	Title       string
	Description string
}

type CompleteIssueInput struct {
	IssueID string
}

type ListIssuesInput struct {
	TeamID string
}

type IssueOutput struct {
	ID          string
	TeamID      string
	Title       string
	Description string
	Status      string
}

type issueUseCase struct {
	issueRepository domain.IssueRepository
}

func NewIssueUseCase(issueRepository domain.IssueRepository) IssueUseCase {
	return &issueUseCase{issueRepository: issueRepository}
}

func (uc *issueUseCase) CreateIssue(ctx context.Context, input CreateIssueInput) (*IssueOutput, error) {
	teamID, err := domain.NewTeamID(input.TeamID)
	if err != nil {
		return nil, ErrInvalidIssueRequest.Wrap(err)
	}

	issue, err := domain.NewIssue(teamID, input.Title, input.Description)
	if err != nil {
		return nil, ErrInvalidIssueRequest.Wrap(err)
	}

	if err := uc.issueRepository.Save(ctx, issue); err != nil {
		return nil, ErrUnexpected.Wrap(err)
	}

	return toOutput(issue), nil
}

func (uc *issueUseCase) CompleteIssue(ctx context.Context, input CompleteIssueInput) (*IssueOutput, error) {
	issueID, err := domain.NewIssueID(input.IssueID)
	if err != nil {
		return nil, ErrInvalidIssueRequest.Wrap(err)
	}

	issue, err := uc.issueRepository.FindByID(ctx, issueID)
	if err != nil {
		return nil, toApplicationError(err)
	}

	if err := issue.Complete(); err != nil {
		return nil, toApplicationError(err)
	}

	if err := uc.issueRepository.Save(ctx, issue); err != nil {
		return nil, ErrUnexpected.Wrap(err)
	}

	return toOutput(issue), nil
}

func (uc *issueUseCase) ListIssues(ctx context.Context, input ListIssuesInput) ([]IssueOutput, error) {
	teamID, err := domain.NewTeamID(input.TeamID)
	if err != nil {
		return nil, ErrInvalidIssueRequest.Wrap(err)
	}

	issues, err := uc.issueRepository.ListByTeam(ctx, teamID)
	if err != nil {
		return nil, ErrUnexpected.Wrap(err)
	}

	outputs := make([]IssueOutput, 0, len(issues))
	for _, issue := range issues {
		outputs = append(outputs, *toOutput(issue))
	}
	return outputs, nil
}

func toApplicationError(err error) error {
	switch {
	case errors.Is(err, domain.ErrIssueNotFound):
		return ErrIssueNotFound.Wrap(err)
	case errors.Is(err, domain.ErrIssueDone):
		return ErrIssueConflict.Wrap(err)
	default:
		return ErrUnexpected.Wrap(err)
	}
}

func toOutput(issue *domain.Issue) *IssueOutput {
	return &IssueOutput{
		ID:          issue.ID().String(),
		TeamID:      issue.TeamID().String(),
		Title:       issue.Title(),
		Description: issue.Description(),
		Status:      string(issue.Status()),
	}
}
