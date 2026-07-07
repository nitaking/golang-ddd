package linear

import (
	"time"

	domain "go-clean-architecture-boilerplate/internal/domain/linear"
)

type issueRecord struct {
	ID          string
	TeamID      string
	Title       string
	Description string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func toRecord(issue *domain.Issue) issueRecord {
	return issueRecord{
		ID:          issue.ID().String(),
		TeamID:      issue.TeamID().String(),
		Title:       issue.Title(),
		Description: issue.Description(),
		Status:      string(issue.Status()),
		CreatedAt:   issue.CreatedAt(),
		UpdatedAt:   issue.UpdatedAt(),
	}
}

func toDomain(record issueRecord) (*domain.Issue, error) {
	issueID, err := domain.NewIssueID(record.ID)
	if err != nil {
		return nil, err
	}

	teamID, err := domain.NewTeamID(record.TeamID)
	if err != nil {
		return nil, err
	}

	return domain.RestoreIssue(
		issueID,
		teamID,
		record.Title,
		record.Description,
		domain.IssueStatus(record.Status),
		record.CreatedAt,
		record.UpdatedAt,
	)
}
