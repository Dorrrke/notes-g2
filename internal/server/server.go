package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Dorrrke/notes-g2/internal"
	notesDomain "github.com/Dorrrke/notes-g2/internal/domain/notes"
	usersDomain "github.com/Dorrrke/notes-g2/internal/domain/users"
	"github.com/Dorrrke/notes-g2/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Repository interface {
	SaveUser(user usersDomain.User) error
	GetUser(login string) (usersDomain.User, error)
	SaveNotes(tasks []notesDomain.Note) error
	GetNotes() ([]notesDomain.Note, error)
	GetNote(nid string) (notesDomain.Note, error)
	Close() error
}

type NotesAPI struct {
	cfg       *internal.Config
	httpServe *http.Server
	repo      Repository
}

func New(cfg *internal.Config, repo Repository) *NotesAPI {
	log := logger.Get()
	log.Debug().Msg("configure Notes API server")
	httpServe := http.Server{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), // 0.0.0.0:8080
	}

	notesAPI := NotesAPI{
		httpServe: &httpServe,
		cfg:       cfg,
		repo:      repo,
	}

	notesAPI.configRoutes()

	return &notesAPI
}

func (nApi *NotesAPI) Run() error {
	log := logger.Get()
	log.Info().Msgf("notes API started on %s", nApi.httpServe.Addr)
	return nApi.httpServe.ListenAndServe()
}

func (nApi *NotesAPI) Stop(ctx context.Context) error {
	return nApi.httpServe.Shutdown(ctx)
}

func (nApi *NotesAPI) configRoutes() {
	log := logger.Get()
	log.Debug().Msg("configure routes")
	router := gin.Default()
	users := router.Group("/users")
	{
		users.GET("/porfile")
		users.POST("/register", nApi.register)
		users.POST("/login", nApi.login)
	}
	notes := router.Group("/notes")
	{
		notes.GET("/", nApi.getTasks)
		notes.GET("/:id", nApi.getTask)
		notes.POST("/", nApi.saveTasks)
		notes.PUT("/:id")
		notes.DELETE("/:id")
	}

	nApi.httpServe.Handler = router
}
