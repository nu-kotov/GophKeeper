package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/nu-kotov/GophKeeper/internal/config"
	"github.com/nu-kotov/GophKeeper/internal/handler"
	"github.com/nu-kotov/GophKeeper/internal/logger"
	"github.com/nu-kotov/GophKeeper/internal/service"
	"github.com/nu-kotov/GophKeeper/internal/storage"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	idleConnsClosed := make(chan struct{})
	sigForShutdown := make(chan os.Signal, 1)
	signal.Notify(sigForShutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	if err := logger.NewLogger("info"); err != nil {
		logger.Log.Info(fmt.Sprintf("Error initialize zap logger: %s", err.Error()))
		return err
	}

	config, err := config.NewConfig()
	if err != nil {
		logger.Log.Info(fmt.Sprintf("Error initialize config: %s", err.Error()))
		return err
	}

	pgStor, err := storage.NewPgStorage(config)
	if err != nil {
		logger.Log.Info(fmt.Sprintf("Error pg connection: %s", err.Error()))
		return err
	}
	miniioStor, err := storage.NewMiniIOStorage(config)
	if err != nil {
		logger.Log.Info(fmt.Sprintf("Error miniio connection: %s", err.Error()))
		return err
	}

	usersStorage := storage.NewUsersStorage(pgStor)
	textStorage := storage.NewTextStorage(pgStor)
	credentialsStorage := storage.NewCredentialsStorage(pgStor)
	cardStorage := storage.NewCardStorage(pgStor)
	binaryStorage := storage.NewBinaryStorage(miniioStor)

	usersService := service.NewUsersService(usersStorage)
	textService := service.NewTextService(textStorage)
	credentialsService := service.NewCredentialsService(credentialsStorage)
	cardService := service.NewCardService(cardStorage)
	binaryService := service.NewBinaryService(binaryStorage)

	r := chi.NewRouter()

	handler.NewUsersHandler(r, config, usersService)
	handler.NewTextHandler(r, config, textService)
	handler.NewCredentialsHandler(r, config, credentialsService)
	handler.NewCardHandler(r, config, cardService)
	handler.NewBinaryHandler(r, config, binaryService)

	defer pgStor.Close()

	server := &http.Server{
		Addr:    config.RunAddr,
		Handler: r,
	}

	go func() error {
		if config.EnableHTTPS {
			manager := &autocert.Manager{
				Cache:      autocert.DirCache("cache-dir"),
				Prompt:     autocert.AcceptTOS,
				HostPolicy: autocert.HostWhitelist("localhost:8181"),
			}
			server.TLSConfig = manager.TLSConfig()
			err = server.ListenAndServeTLS("", "")
			if err != nil && err != http.ErrServerClosed {
				return fmt.Errorf("error ListenAndServeTLS: %w", err)
			}
		} else {
			err = server.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				return fmt.Errorf("error ListenAndServe: %w", err)
			}
		}
		return nil
	}()

	<-sigForShutdown

	logger.Log.Info("shutdown signal received...")

	if err := usersStorage.Stor.Close(); err != nil {
		return fmt.Errorf("error closing users store: %w", err)
	}
	if err := textStorage.Stor.Close(); err != nil {
		return fmt.Errorf("error closing text store: %w", err)
	}
	if err := credentialsStorage.Stor.Close(); err != nil {
		return fmt.Errorf("error closing credentials store: %w", err)
	}
	if err := cardStorage.Stor.Close(); err != nil {
		return fmt.Errorf("error closing card store: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	} else {
		logger.Log.Info("server shutdown gracefully")
	}
	close(idleConnsClosed)

	return nil
}
