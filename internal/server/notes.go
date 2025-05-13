package server

import (
	"net/http"

	notesDomain "github.com/Dorrrke/notes-g2/internal/domain/notes"
	"github.com/gin-gonic/gin"
)

func (s *NotesAPI) saveTasks(ctx *gin.Context) {
	var notes []notesDomain.Note
	if err := ctx.ShouldBindBodyWithJSON(&notes); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := s.repo.SaveNotes(notes); err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.String(http.StatusOK, "notes saved")
}

func (s *NotesAPI) getTasks(ctx *gin.Context) {
	notes, err := s.repo.GetNotes()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	notesResp := []notesDomain.NoteResponseFormat{}
	for _, note := range notes {
		notesResp = append(notesResp, notesDomain.NoteResponse(note))
	}

	ctx.JSON(http.StatusOK, notesResp)
}

func (s *NotesAPI) getTask(ctx *gin.Context) {
	nid := ctx.Param("id")
	note, err := s.repo.GetNote(nid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	noteResp := notesDomain.NoteResponse(note)

	ctx.JSON(http.StatusOK, noteResp)
}
