package internal

import (
	"cmp"
	"os"
)

type Config struct {
	DbDSN string
}

func ReadConfgi() Config {
	var cfg Config

	cfg.DbDSN = cmp.Or(os.Getenv("DB_DSN"), "postgres://user:password@localhost:5432/notes?sslmode=disable")

	return cfg
}
