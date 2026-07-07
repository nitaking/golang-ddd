package linear

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type IssueID string

func NewIssueID(value string) (IssueID, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", ErrEmptyIssueID
	}
	return IssueID(trimmed), nil
}

func NewIssueIDFromUUID() IssueID {
	return IssueID(uuid.NewString())
}

func (id IssueID) String() string {
	return string(id)
}

type TeamID string

func NewTeamID(value string) (TeamID, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", ErrEmptyTeamID
	}
	return TeamID(trimmed), nil
}

func (id TeamID) String() string {
	return string(id)
}

type IssueStatus string

const (
	IssueStatusTodo IssueStatus = "todo"
	IssueStatusDone IssueStatus = "done"
)

type Issue struct {
	id          IssueID
	teamID      TeamID
	title       string
	description string
	status      IssueStatus
	createdAt   time.Time
	updatedAt   time.Time
}

var (
	ErrEmptyIssueID    = errors.New("issue id is empty")
	ErrEmptyTeamID     = errors.New("team id is empty")
	ErrEmptyTitle      = errors.New("issue title is empty")
	ErrIssueDone       = errors.New("issue is already done")
	ErrIssueNotFound   = errors.New("issue not found")
	ErrInvalidArgument = errors.New("invalid argument")
)

func NewIssue(teamID TeamID, title, description string) (*Issue, error) {
	trimmedTitle := strings.TrimSpace(title)
	if trimmedTitle == "" {
		return nil, ErrEmptyTitle
	}

	now := time.Now().UTC()
	return &Issue{
		id:          NewIssueIDFromUUID(),
		teamID:      teamID,
		title:       trimmedTitle,
		description: strings.TrimSpace(description),
		status:      IssueStatusTodo,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

func RestoreIssue(id IssueID, teamID TeamID, title, description string, status IssueStatus, createdAt, updatedAt time.Time) (*Issue, error) {
	trimmedTitle := strings.TrimSpace(title)
	if trimmedTitle == "" {
		return nil, ErrEmptyTitle
	}

	return &Issue{
		id:          id,
		teamID:      teamID,
		title:       trimmedTitle,
		description: strings.TrimSpace(description),
		status:      status,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}, nil
}

func (i *Issue) Complete() error {
	if i.status == IssueStatusDone {
		return ErrIssueDone
	}

	i.status = IssueStatusDone
	i.updatedAt = time.Now().UTC()
	return nil
}

func (i *Issue) ID() IssueID {
	return i.id
}

func (i *Issue) TeamID() TeamID {
	return i.teamID
}

func (i *Issue) Title() string {
	return i.title
}

func (i *Issue) Description() string {
	return i.description
}

func (i *Issue) Status() IssueStatus {
	return i.status
}

func (i *Issue) CreatedAt() time.Time {
	return i.createdAt
}

func (i *Issue) UpdatedAt() time.Time {
	return i.updatedAt
}
