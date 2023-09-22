package main

import (
	"Sber/app/internal/cache"
	"Sber/app/internal/server"
	"Sber/app/internal/task"
	"Sber/app/pkg/config"
	"Sber/app/pkg/logger"
	postgres "Sber/app/pkg/storage"
	"context"
	"errors"
	"flag"
	"github.com/jackc/pgx/v4"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title SberTask
// @host localhost:3003
func main() {
	log := logger.GetLogger()
	log.Info("Logger initialized")

	configPath := flag.String("config-path", "config.yml", "path for application configuration file")
	cfg := config.GetConfig(*configPath, ".env")
	log.Info("Loaded config file")

	dbConn, err := postgres.ConnectDB(*cfg)
	if err != nil {
		log.Error("cannot connect to database", err)
	}
	log.Info("Connected to database")

	allCache := cache.NewCache()
	if err = loadAllCache(log, dbConn, allCache); err != nil {
		log.Error("Failed to preload caches:", err)
	}
	log.Info("Records from the database are added to the cache")

	numItems := len(allCache.Task)
	log.Infof("Cache after loading data contains %d items", numItems)

	router := httprouter.New()
	log.Info("Initialized httprouter")

	srv := server.NewServer(cfg, router, &log, allCache)
	log.Info("Starting the server")

	quit := make(chan os.Signal, 1)
	signals := []os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM}
	signal.Notify(quit, signals...)

	go func() {
		if err = srv.Run(dbConn); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("cannot run the server", err)
		}
	}()
	log.Info("Server has been started ", slog.String("host", cfg.HTTP.Host), slog.String("port", cfg.HTTP.Port))

	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		dbCloseCtx, dbCloseCancel := context.WithTimeout(
			context.Background(),
			time.Duration(cfg.PostgreSQL.ShutdownTimeout)*time.Second,
		)
		defer dbCloseCancel()
		err = dbConn.Close(dbCloseCtx)
		if err != nil {
			log.Error("failed to close database connection:", err)
		}
		log.Info("Closed database connection")
		cancel()
	}()

	if err = srv.Shutdown(ctx); err != nil {
		log.Error("server shutdown failed:", err)
	}
	log.Info("Server has been shutted down")
}

func loadAllCache(log logger.Logger, dbConn *pgx.Conn, cache *cache.Cache) error {
	if err := task.CacheForTask(dbConn, cache); err != nil {
		log.Error("Failed to load task data into cache:", err)
		return err
	}
	return nil
}
