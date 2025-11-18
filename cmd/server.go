package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/gookit/slog"
	"github.com/kstsm/wb-comment-tree/database"
	"github.com/kstsm/wb-comment-tree/internal/config"
	"github.com/kstsm/wb-comment-tree/internal/handler"
	"github.com/kstsm/wb-comment-tree/internal/repository"
	"github.com/kstsm/wb-comment-tree/internal/service"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	cfg := config.GetConfig()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	conn := database.InitPostgres(ctx)
	defer conn.Close()

	repo := repository.NewRepository(conn)
	svc := service.NewService(repo)
	router := handler.NewHandler(svc)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler: router.NewRouter(),
	}

	errChan := make(chan error, 1)

	go func() {
		slog.Infof("Starting server on %s:%d", cfg.Server.Host, cfg.Server.Port)
		errChan <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		slog.Info("Finishing the server...")
	case err := <-errChan:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Fatal("Error starting server", "error", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Error while shutting down the server", "error", err)
	}
}
