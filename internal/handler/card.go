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

// CardStorage - интерфейс хранилища для работы с данными банковских карт.
type CardStorage interface {
	InsertCardData(context.Context, string, *models.CardData) error
	SelectCardData(context.Context, string, string) (string, error)
	DeleteCardData(context.Context, string, string) error
}

// CardHandler - структура хендлера http сервиса для работы с данными банковских карт.
type CardHandler struct {
	Config  *config.Config
	Storage CardStorage
}

// NewCardHandler - конструктор хендлера http сервиса для работы с данными банковских карт.
func NewCardHandler(router *chi.Mux, cfg *config.Config, storage CardStorage) {

	handler := &CardHandler{
		Config:  cfg,
		Storage: storage,
	}

	middlewareStack := middleware.Chain(
		middleware.RequestLogger,
	)

	router.Post(`/api/card/add`, middlewareStack(handler.AddCard()))
	router.Post(`/api/card/get`, middlewareStack(handler.GetCard()))
	router.Post(`/api/card/delete`, middlewareStack(handler.DeleteCard()))

}

// AddCard сохраняет данные банковской карты в pg.
func (handler *CardHandler) AddCard() http.HandlerFunc {
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

		var jsonBody models.CardData
		if err = json.Unmarshal(body, &jsonBody); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		err = handler.Storage.InsertCardData(req.Context(), userID, &jsonBody)
		if err != nil {
			logger.Log.Info(err.Error())
			if errors.Is(err, dberrors.ErrConflict) {
				res.Header().Set("Content-Type", "text/plain")
				res.WriteHeader(http.StatusConflict)
				io.WriteString(res, string("Card "+jsonBody.DataID+" already exists"))
				return
			}
			http.Error(res, "Insert card data error", http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusCreated)
		io.WriteString(res, "Card "+jsonBody.DataID+" added successfully")
	}
}

// GetCard получает данные банковской карты из pg.
func (handler *CardHandler) GetCard() http.HandlerFunc {
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

		var jsonBody models.CardID
		if err = json.Unmarshal(body, &jsonBody); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		card, err := handler.Storage.SelectCardData(req.Context(), userID, jsonBody.DataID)
		if err != nil {
			logger.Log.Info(err.Error())

			if errors.Is(err, dberrors.ErrNotFound) {
				res.Header().Set("Content-Type", "text/plain")
				res.WriteHeader(http.StatusNotFound)
				io.WriteString(res, string("Card "+jsonBody.DataID+" not found"))
				return
			}

			http.Error(res, "Select card data error", http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusOK)
		io.WriteString(res, card)
	}
}

// DeleteCard удаляет данные банковской карты из pg.
func (handler *CardHandler) DeleteCard() http.HandlerFunc {
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

		var jsonBody models.CardID
		if err = json.Unmarshal(body, &jsonBody); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		err = handler.Storage.DeleteCardData(req.Context(), userID, jsonBody.DataID)
		if err != nil {
			logger.Log.Info(err.Error())
			http.Error(res, "Delete card data error", http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusOK)
		io.WriteString(res, "Card "+jsonBody.DataID+" deleted successfully")
	}
}
