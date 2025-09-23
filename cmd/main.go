package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/nu-kotov/GophKeeper/internal/config"
	"github.com/nu-kotov/GophKeeper/internal/handler"
	"github.com/nu-kotov/GophKeeper/internal/logger"
	"github.com/nu-kotov/GophKeeper/internal/storage"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
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

	usersStorage := storage.NewUsersStorage(pgStor)
	textStorage := storage.NewTextStorage(pgStor)

	r := chi.NewRouter()

	handler.NewUsersHandler(r, config, usersStorage)
	handler.NewTextHandler(r, config, textStorage)

	defer pgStor.Close()

	err = http.ListenAndServe(config.RunAddr, r)
	if err != nil {
		logger.Log.Info(fmt.Sprintf("Error starting server: %s", err.Error()))
		return err
	}

	return nil
}
