package application

import (
	"net/http"
	"github.com/troxanna/pr-chat-backend/internal/config"
	"github.com/troxanna/pr-chat-backend/internal/application/rest"
	"context"
	"os/signal"
	"syscall"
	"golang.org/x/sync/errgroup"
	"github.com/troxanna/pr-chat-backend/internal/domain/services"
	"github.com/troxanna/pr-chat-backend/internal/infrastructure/persistence"
	"github.com/jackc/pgx/v5/pgxpool"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net"
	"errors"
	"github.com/troxanna/pr-chat-backend/internal/db"
	"log"
)

type App struct {
	name string
	httpServer *http.Server
	cfg      config.Config
	deferred []func()
	postgresClient *pgxpool.Pool

	dbCompetencyMatrix persistence.DBCompetencyMatrix
	competencyMatrixService *service.CompetencyMatrix

}

func New(name string, cfg config.Config) *App {
	//nolint:exhaustruct
	return &App{
		name:    name,
		cfg:     cfg,
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

	db, _ := db.NewPostgres(ctx, "10.10.169.1")
	app.postgresClient = db.Pool
	log.Println(db)

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
