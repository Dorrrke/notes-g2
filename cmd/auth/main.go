package main

import (
	"context"
	"log"
	"net"

	internal "github.com/Dorrrke/notes-g2/internal/auth"
	"github.com/Dorrrke/notes-g2/internal/auth/gen/grpcapi"
	"github.com/Dorrrke/notes-g2/internal/auth/repository"
	"github.com/Dorrrke/notes-g2/internal/auth/usecase"

	"google.golang.org/grpc"
)

func main() {
	// конфигурация приложения
	cfg := internal.ReadConfgi()

	// инициализация базы данных
	repo, err := repository.NewRepository(context.Background(), cfg.DbDSN)
	if err != nil {
		panic(err)
	}

	usecase := usecase.NewUserUsecase(repo)

	grpcSrv := grpc.NewServer()

	grpcapi.RegisterGrpcAPI(grpcSrv, usecase)

	log.Println("grpc server started...")

	listener, err := net.Listen("tcp", ":9091")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if err := grpcSrv.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
