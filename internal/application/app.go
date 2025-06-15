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
	"github.com/troxanna/pr-chat-backend/internal/application/rest"
	"github.com/troxanna/pr-chat-backend/internal/config"
	"github.com/troxanna/pr-chat-backend/pkg/openai"
	"github.com/troxanna/pr-chat-backend/pkg/bot"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/troxanna/pr-chat-backend/internal/domain/services"
	"github.com/troxanna/pr-chat-backend/internal/infrastructure/integration"
	"github.com/troxanna/pr-chat-backend/internal/infrastructure/persistence"
	"golang.org/x/sync/errgroup"
	"github.com/jmoiron/sqlx"
)


type App struct {
	name           string
	httpServer     *http.Server
	cfg            config.Config
	deferred       []func()
	postgresClient *sqlx.DB
	clientAI       integration.ChatGPTService
	clientBot 		integration.TelegramBotService

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

	postgresClient, err := sqlx.ConnectContext(ctx, "pgx", app.cfg.Postgres.DSN)
	if err != nil {
		return fmt.Errorf("sqlx.ConnectContext: %w", err)
	}

	postgresClient.SetMaxOpenConns(app.cfg.Postgres.MaxOpenConns)
	postgresClient.SetMaxIdleConns(app.cfg.Postgres.MaxIdleConns)
	postgresClient.SetConnMaxLifetime(app.cfg.Postgres.ConnMaxLifetime)

	app.postgresClient = postgresClient

	// gptClient := openai.NewClient(
	// 	&http.Client{Transport: http.DefaultTransport},
	// 	app.cfg.ClientAI.BaseURL,
	// 	app.cfg.ClientAI.APIKey,
	// )
	// gptService := integration.NewChatGPTService(gptClient, "gpt-3.5-turbo")

	// result, err := gptService.AskUser(ctx, messageQuestionTemplate)
	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println(result)

	app.dbCompetencyMatrix = persistence.NewDBCompetencyMatrix(app.postgresClient)
	app.competencyMatrixService = service.NewCompetencyMatrix(app.dbCompetencyMatrix)

	app.runPollingBot(ctx)
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

func (app *App) runPollingBot(ctx context.Context) {
	gptClient := openai.NewClient(
		&http.Client{Transport: http.DefaultTransport},
		app.cfg.ClientAI.BaseURL,
		app.cfg.ClientAI.APIKey,
	)
	app.clientAI = integration.NewChatGPTService(gptClient, "gpt-3.5-turbo")

	botClient, err := bot.NewClient(app.cfg.Telegram.BotToken)
	if err != nil {
		log.Fatalf("bot exited with error: %v", err)
	}
	app.clientBot = integration.NewTelegramBotService(botClient, app.clientAI)
	log.Println("Telegram bot started...")
	go func() {
        if err := app.clientBot.Start(ctx); err != nil {
            log.Printf("Telegram bot error: %v", err)
        }
    }()
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
		Addr:              app.cfg.HTTP.ListenAddressAdmin,
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
