package app

import (
	"context"
	"fmt"
	grpccleint "github.com/shamank/ai-marketplace-api-gateway/internal/clients/stats-service/grpc"
	"github.com/shamank/ai-marketplace-api-gateway/internal/config"
	"github.com/shamank/ai-marketplace-api-gateway/internal/delivery/http"
	"log/slog"
	"os"
)

type App struct {
	cfg        *config.Config
	log        *slog.Logger
	httpServer *HTTPServer
}

func NewApp(cfg *config.Config) *App {

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	grpcAddr := fmt.Sprintf("%s:%d", cfg.StatsService.Host, cfg.StatsService.Port)
	grpcClient, err := grpccleint.NewStatsServiceClient(context.Background(), log, grpcAddr)
	if err != nil {
		log.Error("failed to create grpc client", "error", err)
		panic(err)
	}

	handler := http.NewHandler(grpcClient, log)

	httpAddr := fmt.Sprintf("%s:%d", cfg.HTTPServer.Host, cfg.HTTPServer.Port)

	httpServer := NewHTTPServer(httpAddr, handler.InitAPIRoutes(), cfg.HTTPServer.Timeout)

	return &App{
		cfg:        cfg,
		log:        log,
		httpServer: httpServer,
	}
}

func (app *App) Run() error {
	if err := app.httpServer.Run(); err != nil {
		return err
	}
	return nil
}

func (app *App) Stop(ctx context.Context) error {
	if err := app.httpServer.Stop(ctx); err != nil {
		return err
	}
	return nil
}
