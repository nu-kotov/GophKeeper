package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"github.com/nu-kotov/GophKeeper/internal/auth"
	"github.com/nu-kotov/GophKeeper/internal/config"
	"github.com/nu-kotov/GophKeeper/internal/dberrors"
	"github.com/nu-kotov/GophKeeper/internal/logger"
	"github.com/nu-kotov/GophKeeper/internal/middleware"
	"github.com/nu-kotov/GophKeeper/internal/models"
)

type CredentialsStorage interface {
	InsertCredentialsData(context.Context, string, *models.Credentials) error
	SelectCredentialsData(context.Context, string, string) (*models.Credentials, error)
	DeleteCredentialsData(context.Context, string, string) error
}

type CredentialsHandler struct {
	Config  *config.Config
	Storage CredentialsStorage
}

func NewCredentialsHandler(router *chi.Mux, cfg *config.Config, storage CredentialsStorage) {

	handler := &CredentialsHandler{
		Config:  cfg,
		Storage: storage,
	}

	middlewareStack := middleware.Chain(
		middleware.RequestLogger,
	)

	router.Post(`/api/credentials/add`, middlewareStack(handler.AddCredentials()))
	router.Post(`/api/credentials/get`, middlewareStack(handler.GetCredentials()))
	router.Post(`/api/credentials/delete`, middlewareStack(handler.DeleteCredentials()))

}

func (handler *CredentialsHandler) AddCredentials() http.HandlerFunc {
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

		var jsonBody models.Credentials
		if err = json.Unmarshal(body, &jsonBody); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		err = handler.Storage.InsertCredentialsData(req.Context(), userID, &jsonBody)
		if err != nil {
			logger.Log.Info(err.Error())
			if errors.Is(err, dberrors.ErrConflict) {
				res.Header().Set("Content-Type", "text/plain")
				res.WriteHeader(http.StatusConflict)
				io.WriteString(res, string("Credentials "+jsonBody.DataID+" already exists"))
				return
			}
			http.Error(res, "Insert credentials data error", http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusCreated)
		io.WriteString(res, "Credentials "+jsonBody.DataID+" added successfully")
	}
}

func (handler *CredentialsHandler) GetCredentials() http.HandlerFunc {
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

		var jsonBody models.CredentialsID
		if err = json.Unmarshal(body, &jsonBody); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		credentials, err := handler.Storage.SelectCredentialsData(req.Context(), userID, jsonBody.DataID)
		if err != nil {
			logger.Log.Info(err.Error())

			if errors.Is(err, dberrors.ErrNotFound) {
				res.Header().Set("Content-Type", "text/plain")
				res.WriteHeader(http.StatusNotFound)
				io.WriteString(res, string("Credentials "+jsonBody.DataID+" not found"))
				return
			}

			http.Error(res, "Select credentials data error", http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)

		JSONResp, err := json.Marshal(credentials)
		if err != nil {
			logger.Log.Info("Failed to marshal credentials", zap.Error(err))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = res.Write(JSONResp)

		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (handler *CredentialsHandler) DeleteCredentials() http.HandlerFunc {
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

		var jsonBody models.CredentialsID
		if err = json.Unmarshal(body, &jsonBody); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		err = handler.Storage.DeleteCredentialsData(req.Context(), userID, jsonBody.DataID)
		if err != nil {
			logger.Log.Info(err.Error())
			http.Error(res, "Delete credentials data error", http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusOK)
		io.WriteString(res, "Credentials "+jsonBody.DataID+" deleted successfully")
	}
}
