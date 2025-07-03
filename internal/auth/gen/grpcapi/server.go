package grpcapi

import (
	"context"
	"log"

	authservicev1 "github.com/Dorrrke/notes-g2/internal/auth/gen"
	"github.com/Dorrrke/notes-g2/internal/auth/models"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Usecase interface {
	Register(models.User) (string, error)
	Login(models.UserRequest) (string, error)
}

type GrpcAPI struct {
	authservicev1.UnimplementedAuthServiceServer
	usecase Usecase
}

func RegisterGrpcAPI(gRPC *grpc.Server, usecase Usecase) {
	authservicev1.RegisterAuthServiceServer(gRPC, &GrpcAPI{usecase: usecase})
}

func (gapi *GrpcAPI) Register(ctx context.Context, req *authservicev1.RegisterRequest) (*authservicev1.AuthResponse, error) {
	log.Println("call register method")

	user := models.User{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	uid, err := gapi.usecase.Register(user)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &authservicev1.AuthResponse{
		Token:   uid,
		Message: "success",
	}, nil
}

func (gapi *GrpcAPI) Login(ctx context.Context, req *authservicev1.LoginRequest) (*authservicev1.AuthResponse, error) {
	log.Println("call login method")

	user := models.UserRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	uid, err := gapi.usecase.Login(user)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &authservicev1.AuthResponse{
		Token:   uid,
		Message: "success",
	}, nil
}
