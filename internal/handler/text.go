package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/nu-kotov/GophKeeper/internal/auth"
	"github.com/nu-kotov/GophKeeper/internal/config"
	"github.com/nu-kotov/GophKeeper/internal/keeper_errors"
	"github.com/nu-kotov/GophKeeper/internal/logger"
	"github.com/nu-kotov/GophKeeper/internal/middleware"
	"github.com/nu-kotov/GophKeeper/internal/models"
	"github.com/nu-kotov/GophKeeper/internal/service"
)

// TextHandler - структура хендлера http сервиса для работы с текстом.
type TextHandler struct {
	config  *config.Config
	service *service.TextService
}

// NewTextHandler - конструктор хендлера http сервиса для работы с текстом.
func NewTextHandler(router *chi.Mux, cfg *config.Config, service *service.TextService) {

	handler := &TextHandler{
		config:  cfg,
		service: service,
	}

	middlewareStack := middleware.Chain(
		middleware.RequestLogger,
	)

	router.Post(`/api/text/add`, middlewareStack(handler.AddText()))
	router.Post(`/api/text/get`, middlewareStack(handler.GetText()))
	router.Post(`/api/text/delete`, middlewareStack(handler.DeleteText()))

}

// AddText сохраняет текст в pg.
func (handler *TextHandler) AddText() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		token, err := req.Cookie("token")

		if err != nil {
			logger.Log.Info(err.Error())
			res.WriteHeader(http.StatusUnauthorized)
			io.WriteString(res, string("Please log in"))
			return
		}

		userID, err := auth.GetUserID(token.Value, handler.config.SecretKey)
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

		err = handler.service.AddText(req.Context(), userID, &jsonBody)
		if err != nil {
			logger.Log.Info(err.Error())

			if errors.Is(err, keeper_errors.ErrConflict) {
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

// GetText получает текст из pg.
func (handler *TextHandler) GetText() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		token, err := req.Cookie("token")

		if err != nil {
			logger.Log.Info(err.Error())
			res.WriteHeader(http.StatusUnauthorized)
			io.WriteString(res, string("Please log in"))
			return
		}

		userID, err := auth.GetUserID(token.Value, handler.config.SecretKey)
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

		txt, err := handler.service.GetText(req.Context(), userID, jsonBody.DataID)
		if err != nil {
			logger.Log.Info(err.Error())

			if errors.Is(err, keeper_errors.ErrNotFound) {
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

// DeleteText удаляет текст из pg.
func (handler *TextHandler) DeleteText() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		token, err := req.Cookie("token")

		if err != nil {
			logger.Log.Info(err.Error())
			res.WriteHeader(http.StatusUnauthorized)
			io.WriteString(res, string("Please log in"))
			return
		}

		userID, err := auth.GetUserID(token.Value, handler.config.SecretKey)
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

		err = handler.service.DeleteText(req.Context(), userID, jsonBody.DataID)
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
