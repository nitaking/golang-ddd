package linear

import (
	"errors"
	"net/http"

	"go-clean-architecture-boilerplate/internal/application/application_error"
	application "go-clean-architecture-boilerplate/internal/application/linear"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	issueUseCase application.IssueUseCase
}

func NewController(issueUseCase application.IssueUseCase) *Controller {
	return &Controller{issueUseCase: issueUseCase}
}

func (ctrl *Controller) RegisterRoutes(router gin.IRouter) {
	router.POST("/linear/issues", ctrl.CreateIssue)
	router.GET("/linear/teams/:teamID/issues", ctrl.ListIssues)
	router.POST("/linear/issues/:issueID/complete", ctrl.CompleteIssue)
}

type createIssueRequest struct {
	TeamID      string `json:"teamId" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

func (ctrl *Controller) CreateIssue(c *gin.Context) {
	var req createIssueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.handleError(c, application.ErrInvalidIssueRequest.Wrap(err))
		return
	}

	issue, err := ctrl.issueUseCase.CreateIssue(c.Request.Context(), application.CreateIssueInput{
		TeamID:      req.TeamID,
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		ctrl.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, issue)
}

func (ctrl *Controller) ListIssues(c *gin.Context) {
	issues, err := ctrl.issueUseCase.ListIssues(c.Request.Context(), application.ListIssuesInput{
		TeamID: c.Param("teamID"),
	})
	if err != nil {
		ctrl.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"issues": issues})
}

func (ctrl *Controller) CompleteIssue(c *gin.Context) {
	issue, err := ctrl.issueUseCase.CompleteIssue(c.Request.Context(), application.CompleteIssueInput{
		IssueID: c.Param("issueID"),
	})
	if err != nil {
		ctrl.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, issue)
}

func (ctrl *Controller) handleError(c *gin.Context, err error) {
	var appErr *application_error.AppError
	if !errors.As(err, &appErr) {
		c.JSON(http.StatusInternalServerError, gin.H{"code": application_error.CodeUnexpected, "message": "internal server error"})
		return
	}

	status := http.StatusInternalServerError
	switch appErr.Code() {
	case application_error.CodeInvalidRequest:
		status = http.StatusBadRequest
	case application_error.CodeNotFound:
		status = http.StatusNotFound
	case application_error.CodeConflict:
		status = http.StatusConflict
	}

	c.JSON(status, gin.H{"code": appErr.Code(), "message": appErr.Message()})
}
