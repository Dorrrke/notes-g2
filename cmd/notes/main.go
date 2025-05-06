package main

import (
	"context"

	"github.com/Dorrrke/notes-g2/internal"
	dbstorage "github.com/Dorrrke/notes-g2/internal/infrastructure/db-storage"
	inmemory "github.com/Dorrrke/notes-g2/internal/infrastructure/in-memory"
	"github.com/Dorrrke/notes-g2/internal/server"
	"github.com/Dorrrke/notes-g2/pkg/logger"
)

func main() {
	cfg := internal.ReadConfig()

	log := logger.Get(cfg.Debug)

	log.Info().Msg("service starting")

	var repo server.Repository
	var err error
	repo, err = dbstorage.New(context.Background(), "postgres://user:password@localhost:5432/notes?sslmode=disable")
	if err != nil {
		log.Warn().Err(err).Msg("failed to connected to db. Use in memory storage")
		repo = inmemory.New()
	}

	notesAPI := server.New(cfg, repo)

	if err := notesAPI.Run(); err != nil {
		log.Error().Err(err).Msg("failed running server")
	}
}
