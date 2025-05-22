package main

import (
	"context"

	"github.com/Dorrrke/notes-g2/internal"
	dbstorage "github.com/Dorrrke/notes-g2/internal/infrastructure/db-storage"
	inmemory "github.com/Dorrrke/notes-g2/internal/infrastructure/in-memory"
	"github.com/Dorrrke/notes-g2/internal/server"
	"github.com/Dorrrke/notes-g2/pkg/logger"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg, err := internal.ReadConfig()
	if err != nil {
		panic(err)
	}

	log := logger.Get(cfg.Debug)

	log.Info().Msg("service starting")

	var repo server.Repository
	repo, err = dbstorage.New(context.Background(), cfg.DBConnStr)
	if err != nil {
		log.Warn().Err(err).Msg("failed to connected to db. Use in memory storage")
		repo = inmemory.New()
	}
	if err = dbstorage.AppyMigrations(cfg.DBConnStr); err != nil {
		log.Warn().Err(err).Msg("failed to apply migrations. Use in memory storage")
		repo.Close()
		repo = inmemory.New()
	}

	log.Info().Msg("connected to db successfully")

	notesAPI := server.New(cfg, repo)

	if err := notesAPI.Run(); err != nil {
		log.Error().Err(err).Msg("failed running server")
	}
}
