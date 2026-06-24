package di

import (
	application "go-clean-architecture-boilerplate/internal/application/linear"
	"go-clean-architecture-boilerplate/internal/domain/linear"
	memorylinear "go-clean-architecture-boilerplate/internal/infrastructure/memory/linear"
	httplinear "go-clean-architecture-boilerplate/internal/presentation/http/linear"
)

type Container struct {
	IssueRepository linear.IssueRepository
	IssueUseCase    application.IssueUseCase
	IssueController *httplinear.Controller
}

func NewContainer() *Container {
	issueRepository := memorylinear.NewIssueRepository()
	issueUseCase := application.NewIssueUseCase(issueRepository)
	issueController := httplinear.NewController(issueUseCase)

	return &Container{
		IssueRepository: issueRepository,
		IssueUseCase:    issueUseCase,
		IssueController: issueController,
	}
}
