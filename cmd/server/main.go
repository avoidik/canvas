// Package main is the entry point to the server. It reads configuration, sets up logging and error handling,
// handles signals from the OS, and starts and stops the server.
package main

import (
	"canvas/server"
	"canvas/storage"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/maragudk/env"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var release string

func main() {
	os.Exit(start())
}

func start() int {
	_ = env.Load()

	logEnv := env.GetStringOrDefault("LOG_ENV", "development")
	log, err := createLogger(logEnv)
	if err != nil {
		fmt.Println("Error setting up the logger:", err)
		return 1
	}

	log = log.With(zap.String("release", release))

	defer func() {
		_ = log.Sync()
	}()

	host := env.GetStringOrDefault("HOST", "localhost")
	port := env.GetIntOrDefault("PORT", 8081)

	s := server.New(server.Options{
		DB:   createDatabase(log),
		Host: host,
		Port: port,
		Log:  log,
	})

	var errGroup errgroup.Group
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	errGroup.Go(func() error {
		<-ctx.Done()
		if err := s.Stop(); err != nil {
			log.Info("Error stopping server", zap.Error(err))
			return err
		}
		return nil
	})

	if err := s.Start(); err != nil {
		log.Info("Error starting server", zap.Error(err))
		return 1
	}

	if errGroup.Wait(); err != nil {
		return 1
	}
	return 0
}

func createDatabase(log *zap.Logger) *storage.Database {
	return storage.NewDatabase(storage.NewDatabaseOptions{
		Host:                  env.GetStringOrDefault("DB_HOST", "localhost"),
		Port:                  env.GetIntOrDefault("DB_PORT", 5432),
		User:                  env.GetStringOrDefault("DB_USER", ""),
		Password:              env.GetStringOrDefault("DB_PASSWORD", ""),
		Name:                  env.GetStringOrDefault("DB_NAME", ""),
		MaxOpenConnections:    env.GetIntOrDefault("DB_MAX_OPEN_CONNECTIONS", 10),
		MaxIdleConnections:    env.GetIntOrDefault("DB_MAX_IDLE_CONNECTIONS", 10),
		ConnectionMaxLifetime: env.GetDurationOrDefault("DB_CONNECTION_MAX_LIFETIME", time.Hour),
		ConnectionMaxIdleTime: env.GetDurationOrDefault("DB_CONNECTION_MAX_IDLETIME", time.Minute),
		Log:                   log,
	})
}

func createLogger(env string) (*zap.Logger, error) {
	switch env {
	case "production":
		return zap.NewProduction()
	case "development":
		return zap.NewDevelopment()
	default:
		return zap.NewNop(), nil
	}
}

func getStringOrDefault(name, defValue string) string {
	v, ok := os.LookupEnv(name)
	if !ok {
		return defValue
	}
	return v
}

func getIntOrDefault(name string, defValue int) int {
	v, ok := os.LookupEnv(name)
	if !ok {
		return defValue
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return defValue
	}
	return i
}
