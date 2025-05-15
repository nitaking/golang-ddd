package main

import (
	"github.com/gin-gonic/gin"
	"go-clean-architecture-boilerplate/internal/domain/note"
	"go-clean-architecture-boilerplate/internal/infrastructure/memory"
	"go-clean-architecture-boilerplate/internal/presentation/http"
	usecase "go-clean-architecture-boilerplate/internal/usecase/note"
	"log"
)

func main() {
	// 0: memory data
	noteMemory := make(map[note.NoteID]*note.Note)

	// 1. infra: リポジトリ実装を生成
	repository := memory.NewNoteRepository(noteMemory)
	queryRepository := memory.NewNoteQueryRepository(noteMemory)

	//service := memory.NewNoteService(repository)

	// 2. usecase: NoteUseCase を生成
	uc := usecase.NewNoteUseCase(repository, queryRepository)

	// 3. presentation: Controller/Handler を生成
	controller := http.NewNoteController(uc)

	// 4. HTTP サーバを組み立てて起動
	r := gin.Default()
	r.POST("/notes", controller.Create)
	r.GET("/notes", controller.Search)
	log.Println("Listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
