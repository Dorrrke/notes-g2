package server

import (
	"context"
	"net/http"

	usersDomain "github.com/Dorrrke/notes-g2/internal/tracker/domain/users"
	authservicev1 "github.com/Dorrrke/notes-g2/internal/tracker/grpclient"
	"github.com/gin-gonic/gin"
)

func (s *NotesAPI) login(c *gin.Context) {
	var uReq usersDomain.UserRequest

	if err := c.ShouldBindJSON(&uReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resp, err := s.auth.Login(context.Background(), &authservicev1.LoginRequest{
		Email:    uReq.Email,
		Password: uReq.Password,
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, "user logined: %s; login info: %s", resp.GetMessage(), resp.GetToken())
}

func (s *NotesAPI) register(ctx *gin.Context) {
	var uReq usersDomain.User

	if err := ctx.ShouldBindJSON(&uReq); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resp, err := s.auth.Register(context.Background(), &authservicev1.RegisterRequest{
		Name:     uReq.Name,
		Email:    uReq.Email,
		Password: uReq.Password,
	})

	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	// userService := usersService.New(s.repo)

	// userID, err := userService.RegisterUser(uReq)
	// if err != nil {
	// 	ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	// 	return
	// }

	ctx.String(http.StatusOK, "user registered: %s; login info: %s", resp.GetMessage(), resp.GetToken())
}
