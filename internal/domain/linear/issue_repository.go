package linear

import "context"

type IssueRepository interface {
	Save(ctx context.Context, issue *Issue) error
	FindByID(ctx context.Context, id IssueID) (*Issue, error)
	ListByTeam(ctx context.Context, teamID TeamID) ([]*Issue, error)
}
