package main

import (
	"log"
	"os"

	"go-clean-architecture-boilerplate/internal/infrastructure/di"
	presentationhttp "go-clean-architecture-boilerplate/internal/presentation/http"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	container := di.NewContainer()
	router := presentationhttp.NewRouter(container)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
