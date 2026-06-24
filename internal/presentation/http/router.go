package http

import (
	"go-clean-architecture-boilerplate/internal/infrastructure/di"

	"github.com/gin-gonic/gin"
)

func NewRouter(container *di.Container) *gin.Engine {
	router := gin.Default()
	container.IssueController.RegisterRoutes(router)
	return router
}
