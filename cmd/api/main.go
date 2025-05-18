package main

import (
	"github.com/gin-gonic/gin"
	"go-clean-architecture-boilerplate/internal/infrastructure/bolt"
	"go-clean-architecture-boilerplate/internal/presentation/http"
	usecase "go-clean-architecture-boilerplate/internal/usecase/note"
	"log"
)

func main() {
	// Init Bolt DB
	db := bolt.NewBboltDB()
	defer db.Close()

	// 1. infra: リポジトリ実装を生成
	noteRepository := bolt.NewNoteRepository(db)
	queryRepository := bolt.NewNoteQueryRepository(db)

	// 2. usecase: NoteUseCase を生成
	uc := usecase.NewNoteUseCase(noteRepository, queryRepository)

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
