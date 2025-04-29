package main

import (
	"fmt"

	"github.com/Dorrrke/notes-g2/internal"
	inmemory "github.com/Dorrrke/notes-g2/internal/infrastructure/in-memory"
	"github.com/Dorrrke/notes-g2/internal/server"
)

func main() {
	cfg := internal.ReadConfig()
	fmt.Printf("Host: %s\nPort: %d\n", cfg.Host, cfg.Port)

	inMemoryRepo := inmemory.New()

	notesAPI := server.New(cfg, inMemoryRepo)

	notesAPI.Run()
}
