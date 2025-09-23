package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/nu-kotov/GophKeeper/internal/auth"
	"github.com/nu-kotov/GophKeeper/internal/config"
	"github.com/nu-kotov/GophKeeper/internal/logger"
	"github.com/nu-kotov/GophKeeper/internal/middleware"
	"github.com/nu-kotov/GophKeeper/internal/models"
	"github.com/nu-kotov/GophKeeper/internal/storage/postgres"
)

type TextStorage interface {
	InsertTextData(context.Context, string, *models.TextData) error
	SelectTextData(context.Context, string, string) (string, error)
	DeleteTextData(context.Context, string, string) error
}

type TextHandler struct {
	Config  *config.Config
	Storage TextStorage
}

func NewTextHandler(router *chi.Mux, cfg *config.Config, storage TextStorage) {

	handler := &TextHandler{
		Config:  cfg,
		Storage: storage,
	}

	middlewareStack := middleware.Chain(
		middleware.RequestLogger,
	)

	router.Post(`/api/text/add`, middlewareStack(handler.AddText()))
	router.Post(`/api/text/get`, middlewareStack(handler.GetText()))
	router.Post(`/api/text/delete`, middlewareStack(handler.DeleteText()))

}

func (handler *TextHandler) AddText() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		token, err := req.Cookie("token")

		if err != nil {
			logger.Log.Info(err.Error())
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		userID, err := auth.GetUserID(token.Value, handler.Config.SecretKey)
		if err != nil {
			logger.Log.Info(err.Error())
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(req.Body)
		if err != nil {
			logger.Log.Info(err.Error())
			http.Error(res, "Invalid body", http.StatusBadRequest)
			return
		}

		var jsonBody models.TextData
		if err = json.Unmarshal(body, &jsonBody); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		err = handler.Storage.InsertTextData(req.Context(), userID, &jsonBody)
		if err != nil {
			logger.Log.Info(err.Error())
			if errors.Is(err, postgres.ErrConflict) {
				res.Header().Set("Content-Type", "text/plain")
				res.WriteHeader(http.StatusConflict)
				io.WriteString(res, string("Text "+jsonBody.DataID+" already exists"))
				return
			}
			http.Error(res, "Insert text data error", http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusCreated)
		io.WriteString(res, "Текст успешно сохранен")
	}
}

func (handler *TextHandler) GetText() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {}
}

func (handler *TextHandler) DeleteText() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {}
}
