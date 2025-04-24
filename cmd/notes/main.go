package main

import (
	"fmt"

	"github.com/Dorrrke/notes-g2/internal/server"
)

func main() {

	server := server.New("0.0.0.0", "8080")
	server.Run()
	//TODO: Сделать конфигурацию
	panic(fmt.Errorf("not implemented"))
}
