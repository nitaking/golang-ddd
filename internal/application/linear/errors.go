package linear

import "go-clean-architecture-boilerplate/internal/application/application_error"

var (
	ErrInvalidIssueRequest = application_error.New(application_error.CodeInvalidRequest, "invalid issue request")
	ErrIssueNotFound       = application_error.New(application_error.CodeNotFound, "issue not found")
	ErrIssueConflict       = application_error.New(application_error.CodeConflict, "issue state conflict")
	ErrUnexpected          = application_error.New(application_error.CodeUnexpected, "unexpected error")
)
