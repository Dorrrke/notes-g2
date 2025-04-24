package main

import (
	"fmt"

	"github.com/Dorrrke/notes-g2/internal"
	"github.com/Dorrrke/notes-g2/internal/server"
)

func main() {
	cfg := internal.ReadConfig()
	fmt.Printf("Host: %s\nPort: %d\n", cfg.Host, cfg.Port)

	server := server.New(cfg)
	server.Run()
}
