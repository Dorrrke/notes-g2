package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	internal "github.com/Dorrrke/notes-g2/internal/tracker"
	notesDomain "github.com/Dorrrke/notes-g2/internal/tracker/domain/notes"
	usersDomain "github.com/Dorrrke/notes-g2/internal/tracker/domain/users"
	authservicev1 "github.com/Dorrrke/notes-g2/internal/tracker/grpclient"
	"github.com/Dorrrke/notes-g2/pkg/logger"
	"github.com/gin-gonic/gin"
)

const readHeaderTimeout = 5 * time.Second

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
	auth      authservicev1.AuthServiceClient
}

func New(cfg *internal.Config, repo Repository, authClient authservicev1.AuthServiceClient) *NotesAPI {
	log := logger.Get()
	log.Debug().Msg("configure Notes API server")
	notesAPI := NotesAPI{
		cfg:  cfg,
		repo: repo,
		auth: authClient,
	}

	httpServe := http.Server{
		Addr:              fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), // 0.0.0.0:8080
		ReadHeaderTimeout: readHeaderTimeout,
	}

	notesAPI.httpServe = &httpServe

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
	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Notes API")
	})
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
