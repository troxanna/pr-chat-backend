package application

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/troxanna/pr-chat-backend/internal/application/rest"
	"github.com/troxanna/pr-chat-backend/internal/config"
	"github.com/troxanna/pr-chat-backend/internal/db"
	"github.com/troxanna/pr-chat-backend/internal/domain/services"
	"github.com/troxanna/pr-chat-backend/internal/infrastructure/integration"
	"github.com/troxanna/pr-chat-backend/internal/infrastructure/persistence"
	"golang.org/x/sync/errgroup"
	"github.com/google/uuid"
)

type App struct {
	name           string
	httpServer     *http.Server
	cfg            config.Config
	deferred       []func()
	postgresClient *pgxpool.Pool
	clientAI       integration.Client

	dbCompetencyMatrix      persistence.DBCompetencyMatrix
	competencyMatrixService *service.CompetencyMatrix
}

func New(name string, cfg config.Config) *App {
	//nolint:exhaustruct
	return &App{
		name: name,
		cfg:  cfg,
	}
}

func (app *App) Run() error {
	defer app.shutdown()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	defer stop()

	g, ctx := errgroup.WithContext(ctx)

	db, err := db.NewPostgres(ctx, "host=pg-db port=5432 user=postgres password=postgres dbname=app_db sslmode=disable")
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	app.postgresClient = db.Pool
	log.Println(db)

	app.clientAI = integration.NewClient(
		&http.Client{Transport: http.DefaultTransport},
		"app.cfg.ClientAI.BaseURL",
		"OrVrQoQ6T43vk0McGmHOsdvvTiX446RJ",
	)

	app.dbCompetencyMatrix = persistence.NewDBCompetencyMatrix(app.postgresClient)
	log.Println(app.dbCompetencyMatrix)
	app.competencyMatrixService = service.NewCompetencyMatrix(app.dbCompetencyMatrix)
	log.Println(app.competencyMatrixService)

	app.runHTTPServer(ctx, g)

	if err := g.Wait(); err != nil {
		return fmt.Errorf("g.Wait: %w", err)
	}

	return nil
}

func (app *App) shutdown() {
	for _, fn := range app.deferred {
		fn()
	}
}

func (app *App) runHTTPServer(ctx context.Context, g *errgroup.Group) {

	app.httpServer = app.newHTTPServer(ctx)
	uid := uuid.NewString()
	result := false
	log.Println(uid)
	app.clientAI.SendPromptForQuestion(uid)
	for !result {
		result = app.clientAI.GetResultForQuestionRequest(uid)
	} 
	app.clientAI.CleanContextRequest(uid)

	g.Go(func() error {
		go func() {
			<-ctx.Done()

			ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), app.cfg.HTTP.ShutdownTimeout) //nolint:govet
			defer cancel()

			if err := app.httpServer.Shutdown(ctx); err != nil {
				fmt.Errorf("Shutdown error: %w", err)
			}
		}()

		// contextx.LoggerFromContextOrDefault(ctx).Info("http server started", "address", app.cfg.HTTP.ListenAddress)

		if err := app.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("httpServer.ListenAndServe: %w", err)
		}

		// contextx.LoggerFromContextOrDefault(ctx).Info("http server stopped")

		return nil
	})
}

func (app *App) newHTTPServer(ctx context.Context) *http.Server {
	router := chi.NewRouter()

	admServer := rest.NewServerAdmin(app.competencyMatrixService)

	rest.RegisterRoutes(router, admServer)

	return &http.Server{ //nolint:exhaustruct
		Addr:              app.cfg.HTTP.ListenAddressPrivate,
		WriteTimeout:      app.cfg.HTTP.WriteTimeout,
		ReadTimeout:       app.cfg.HTTP.ReadTimeout,
		ReadHeaderTimeout: app.cfg.HTTP.ReadTimeout,
		IdleTimeout:       app.cfg.HTTP.IdleTimeout,
		Handler:           router,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}
}
