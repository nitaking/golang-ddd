package linear

import (
	"errors"
	"testing"
)

func TestNewIssue(t *testing.T) {
	// Arrange
	teamID, err := NewTeamID("team-1")
	if err != nil {
		t.Fatalf("NewTeamID() error = %v", err)
	}

	// Act
	issue, err := NewIssue(teamID, " Fix login ", " description ")

	// Assert
	if err != nil {
		t.Fatalf("NewIssue() error = %v", err)
	}
	if issue.Title() != "Fix login" {
		t.Fatalf("Title() = %q, want %q", issue.Title(), "Fix login")
	}
	if issue.Status() != IssueStatusTodo {
		t.Fatalf("Status() = %q, want %q", issue.Status(), IssueStatusTodo)
	}
}

func TestNewIssueWithEmptyTitle(t *testing.T) {
	// Arrange
	teamID, err := NewTeamID("team-1")
	if err != nil {
		t.Fatalf("NewTeamID() error = %v", err)
	}

	// Act
	_, err = NewIssue(teamID, " ", "description")

	// Assert
	if !errors.Is(err, ErrEmptyTitle) {
		t.Fatalf("NewIssue() error = %v, want %v", err, ErrEmptyTitle)
	}
}

func TestIssueComplete(t *testing.T) {
	// Arrange
	teamID, err := NewTeamID("team-1")
	if err != nil {
		t.Fatalf("NewTeamID() error = %v", err)
	}
	issue, err := NewIssue(teamID, "Fix login", "description")
	if err != nil {
		t.Fatalf("NewIssue() error = %v", err)
	}

	// Act
	err = issue.Complete()

	// Assert
	if err != nil {
		t.Fatalf("Complete() error = %v", err)
	}
	if issue.Status() != IssueStatusDone {
		t.Fatalf("Status() = %q, want %q", issue.Status(), IssueStatusDone)
	}
}
