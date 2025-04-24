package server

import (
	"net/http"

	usersDomain "github.com/Dorrrke/notes-g2/internal/domain/users"
	"github.com/gin-gonic/gin"
)

func (s *Server) login(ctx *gin.Context) {
	var uReq usersDomain.UserRequest

	if err := ctx.ShouldBindJSON(&uReq); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// TODO: вызов бизнес логики авторизации

	ctx.JSON(http.StatusOK, nil)
}
