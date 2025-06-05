package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Dorrrke/notes-g2/internal"
	dbstorage "github.com/Dorrrke/notes-g2/internal/infrastructure/db-storage"
	inmemory "github.com/Dorrrke/notes-g2/internal/infrastructure/in-memory"
	"github.com/Dorrrke/notes-g2/internal/server"
	"github.com/Dorrrke/notes-g2/pkg/logger"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"golang.org/x/sync/errgroup"
)

func gracefulShutdown(cancel context.CancelFunc) {
	log := logger.Get()

	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)

	sig := <-c

	log.Info().Msgf("graceful shutdown with signal: %s", sig)
	cancel()
}

func main() {
	cfg, err := internal.ReadConfig()
	if err != nil {
		panic(err)
	}

	log := logger.Get(cfg.Debug)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go gracefulShutdown(cancel)

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

	group, gCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return notesAPI.Run()
	})

	group.Go(func() error {
		<-gCtx.Done()
		if err := notesAPI.Stop(gCtx); err != nil {
			return err
		}
		if err := repo.Close(); err != nil {
			return err
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("service stopped with error")
		}
	}
	log.Info().Msg("service stopped gracefully")
}
