package http

import (
	"github.com/gin-gonic/gin"
	"go-clean-architecture-boilerplate/internal/domain/note"
	usecase "go-clean-architecture-boilerplate/internal/usecase/note"
	"net/http"
)

type NoteController struct {
	NoteUseCase usecase.NoteUseCase
}

func NewNoteController(uc usecase.NoteUseCase) *NoteController {
	return &NoteController{
		NoteUseCase: uc,
	}
}

func (nc *NoteController) Create(c *gin.Context) {
	var req struct{ Title, Content string }
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := nc.NoteUseCase.CreateNote(c, usecase.CreateNoteInput{
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (nc *NoteController) Search(c *gin.Context) {
	kw := c.Query("q")
	list, err := nc.NoteUseCase.SearchNote(c, kw)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (nc *NoteController) Edit(c *gin.Context) {
	inputID := c.Param("id")
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := nc.NoteUseCase.EditNote(c, usecase.EditNoteInput{
		ID:      note.NoteID(inputID),
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (nc *NoteController) Delete(c *gin.Context) {
	inputID := c.Param("id")
	result, err := nc.NoteUseCase.DeleteNote(c, usecase.DeleteNoteInput{
		ID: note.NoteID(inputID),
	})
	if result == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "note not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": result.ID, "deleted": true})
}
