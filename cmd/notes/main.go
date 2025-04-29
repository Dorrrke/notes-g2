package main

import (
	"github.com/Dorrrke/notes-g2/internal"
	inmemory "github.com/Dorrrke/notes-g2/internal/infrastructure/in-memory"
	"github.com/Dorrrke/notes-g2/internal/server"
	"github.com/Dorrrke/notes-g2/pkg/logger"
)

func main() {
	cfg := internal.ReadConfig()

	log := logger.Get(cfg.Debug)

	log.Info().Msg("service starting")

	inMemoryRepo := inmemory.New()

	notesAPI := server.New(cfg, inMemoryRepo)

	if err := notesAPI.Run(); err != nil {
		log.Error().Err(err).Msg("failed running server")
	}
}
