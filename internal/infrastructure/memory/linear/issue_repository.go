package linear

import (
	"context"
	"sync"

	domain "go-clean-architecture-boilerplate/internal/domain/linear"
)

type IssueRepository struct {
	mu     sync.RWMutex
	issues map[string]issueRecord
}

func NewIssueRepository() *IssueRepository {
	return &IssueRepository{issues: make(map[string]issueRecord)}
}

func (r *IssueRepository) Save(ctx context.Context, issue *domain.Issue) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.issues[issue.ID().String()] = toRecord(issue)
	return nil
}

func (r *IssueRepository) FindByID(ctx context.Context, id domain.IssueID) (*domain.Issue, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	record, ok := r.issues[id.String()]
	r.mu.RUnlock()
	if !ok {
		return nil, domain.ErrIssueNotFound
	}

	return toDomain(record)
}

func (r *IssueRepository) ListByTeam(ctx context.Context, teamID domain.TeamID) ([]*domain.Issue, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	records := make([]issueRecord, 0, len(r.issues))
	for _, record := range r.issues {
		if record.TeamID == teamID.String() {
			records = append(records, record)
		}
	}
	r.mu.RUnlock()

	issues := make([]*domain.Issue, 0, len(records))
	for _, record := range records {
		issue, err := toDomain(record)
		if err != nil {
			return nil, err
		}
		issues = append(issues, issue)
	}

	return issues, nil
}
