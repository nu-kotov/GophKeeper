package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/nu-kotov/GophKeeper/internal/auth"

	"github.com/nu-kotov/GophKeeper/internal/config"
	"github.com/nu-kotov/GophKeeper/internal/dberrors"
	"github.com/nu-kotov/GophKeeper/internal/logger"
	"github.com/nu-kotov/GophKeeper/internal/middleware"
	"github.com/nu-kotov/GophKeeper/internal/models"
)

// UsersStorage - интерфейс хранилища для работы с пользователями.
type UsersStorage interface {
	InsertUserData(context.Context, *models.UserData) error
	SelectUserData(context.Context, *models.UserData) (*models.UserData, error)
}

// UsersHandler - структура хендлера http сервиса для работы с пользователями.
type UsersHandler struct {
	Config  *config.Config
	Storage UsersStorage
}

// NewUsersHandler - конструктор хендлера http сервиса для работы с пользователями.
func NewUsersHandler(router *chi.Mux, cfg *config.Config, storage UsersStorage) {

	handler := &UsersHandler{
		Config:  cfg,
		Storage: storage,
	}

	middlewareStack := middleware.Chain(
		middleware.RequestLogger,
	)

	router.Post(`/api/user/register`, middlewareStack(handler.RegisterUser()))
	router.Post(`/api/user/login`, middlewareStack(handler.LoginUser()))

}

// RegisterUser - метод для регистрации нового пользователя.
func (handler *UsersHandler) RegisterUser() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			logger.Log.Info(err.Error())
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		var jsonBody models.UserData
		if err = json.Unmarshal(body, &jsonBody); err != nil {
			logger.Log.Info(err.Error())
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		passwordHash, err := argon2id.CreateHash(jsonBody.Password, argon2id.DefaultParams)
		if err != nil {
			logger.Log.Info(err.Error())
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonBody.Password = passwordHash
		jsonBody.UserID = uuid.New().String()

		err = handler.Storage.InsertUserData(req.Context(), &jsonBody)
		if err != nil {
			logger.Log.Info(err.Error())
			if errors.Is(err, dberrors.ErrConflict) {
				res.Header().Set("Content-Type", "text/plain")
				res.WriteHeader(http.StatusConflict)
				io.WriteString(res, string("User "+jsonBody.Login+" already exists"))
				return
			}
			http.Error(res, "Register user error", http.StatusInternalServerError)
			return
		}

		value, err := auth.BuildJWTString(
			jsonBody.UserID,
			jsonBody.Login,
			handler.Config.TokenExp,
			handler.Config.SecretKey,
		)
		if err != nil {
			logger.Log.Info(err.Error())
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		cookie := &http.Cookie{
			Name:     "token",
			Value:    value,
			HttpOnly: true,
			Path:     "/",
		}

		http.SetCookie(res, cookie)
		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusOK)
		io.WriteString(res, fmt.Sprintf("User %s registered successfully", jsonBody.Login))
	}
}

// LoginUser - метод для авторизации пользователя.
func (handler *UsersHandler) LoginUser() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			logger.Log.Info(err.Error())
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		var jsonBody models.UserData
		if err = json.Unmarshal(body, &jsonBody); err != nil {
			logger.Log.Info(err.Error())
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		userData, err := handler.Storage.SelectUserData(req.Context(), &jsonBody)
		if err != nil {
			logger.Log.Info(err.Error())
			http.Error(res, "Get user password error", http.StatusInternalServerError)
			return
		}

		match, err := argon2id.ComparePasswordAndHash(jsonBody.Password, userData.Password)
		if err != nil {
			logger.Log.Info(err.Error())
			http.Error(res, "Password comparing error", http.StatusInternalServerError)
			return
		}
		if !match {
			res.Header().Set("Content-Type", "text/plain")
			res.WriteHeader(http.StatusUnauthorized)
			io.WriteString(res, "Uncorrect passwort")
			return
		}

		value, err := auth.BuildJWTString(
			userData.UserID,
			userData.Login,
			handler.Config.TokenExp,
			handler.Config.SecretKey,
		)
		if err != nil {
			logger.Log.Info(err.Error())
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		cookie := &http.Cookie{
			Name:     "token",
			Value:    value,
			HttpOnly: true,
			Path:     "/",
		}

		http.SetCookie(res, cookie)
		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusOK)
		io.WriteString(res, "User "+jsonBody.Login+" authorized")
	}
}
