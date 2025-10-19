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
	"github.com/nu-kotov/GophKeeper/internal/dberrors"
	"github.com/nu-kotov/GophKeeper/internal/logger"
	"github.com/nu-kotov/GophKeeper/internal/middleware"
	"github.com/nu-kotov/GophKeeper/internal/models"
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
			io.WriteString(res, string("Please log in"))
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

			if errors.Is(err, dberrors.ErrConflict) {
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
		io.WriteString(res, "Text "+jsonBody.DataID+" added successfully")
	}
}

func (handler *TextHandler) GetText() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		token, err := req.Cookie("token")

		if err != nil {
			logger.Log.Info(err.Error())
			res.WriteHeader(http.StatusUnauthorized)
			io.WriteString(res, string("Please log in"))
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

		var jsonBody models.TextID
		if err = json.Unmarshal(body, &jsonBody); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		txt, err := handler.Storage.SelectTextData(req.Context(), userID, jsonBody.DataID)
		if err != nil {
			logger.Log.Info(err.Error())

			if errors.Is(err, dberrors.ErrNotFound) {
				res.Header().Set("Content-Type", "text/plain")
				res.WriteHeader(http.StatusNotFound)
				io.WriteString(res, string("Text "+jsonBody.DataID+" not found"))
				return
			}

			http.Error(res, "Select text data error", http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusOK)
		io.WriteString(res, txt)
	}
}

func (handler *TextHandler) DeleteText() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		token, err := req.Cookie("token")

		if err != nil {
			logger.Log.Info(err.Error())
			res.WriteHeader(http.StatusUnauthorized)
			io.WriteString(res, string("Please log in"))
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

		var jsonBody models.TextID
		if err = json.Unmarshal(body, &jsonBody); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		err = handler.Storage.DeleteTextData(req.Context(), userID, jsonBody.DataID)
		if err != nil {
			logger.Log.Info(err.Error())
			http.Error(res, "Delete text data error", http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusOK)
		io.WriteString(res, "Text "+jsonBody.DataID+" deleted successfully")
	}
}
