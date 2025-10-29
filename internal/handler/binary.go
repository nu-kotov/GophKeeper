package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/nu-kotov/GophKeeper/internal/auth"
	"github.com/nu-kotov/GophKeeper/internal/config"
	"github.com/nu-kotov/GophKeeper/internal/logger"
	"github.com/nu-kotov/GophKeeper/internal/middleware"
	"github.com/nu-kotov/GophKeeper/internal/service"
)

// BinaryHandler - структура хендлера http сервиса для работы с бинарными данными.
type BinaryHandler struct {
	config  *config.Config
	service *service.BinaryService
}

// NewBinaryHandler - конструктор хендлера http сервиса для работы с бинарными данными.
func NewBinaryHandler(router *chi.Mux, cfg *config.Config, service *service.BinaryService) {

	handler := &BinaryHandler{
		config:  cfg,
		service: service,
	}

	middlewareStack := middleware.Chain(
		middleware.RequestLogger,
	)

	router.Post(`/api/binary/add`, middlewareStack(handler.AddBinary()))
	router.Post(`/api/binary/get`, middlewareStack(handler.GetBinary()))
	router.Post(`/api/binary/delete`, middlewareStack(handler.DeleteBinary()))

}

// AddBinary сохраняет бинарные данные в miniio.
func (handler *BinaryHandler) AddBinary() http.HandlerFunc {
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

		filename := req.Header.Get("X-Filename")
		if filename == "" {
			http.Error(res, "Missing X-Filename header", http.StatusBadRequest)
			return
		}

		miniioFilename := fmt.Sprintf("%s/%s", userID, filename)

		_, err = handler.service.AddBinary(
			req.Context(),
			handler.config.MiniIOConnection.BucketName,
			miniioFilename,
			req.Body,
		)
		if err != nil {
			logger.Log.Info(err.Error())
			http.Error(res, "File uploaded error", http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusCreated)
		io.WriteString(res, "File "+filename+" uploaded successfully")
	}
}

// GetBinary получает бинарные данные из miniio.
func (handler *BinaryHandler) GetBinary() http.HandlerFunc {
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

		filename := string(body)
		miniioFilename := fmt.Sprintf("%s/%s", userID, filename)

		file, err := handler.service.GetBinary(req.Context(), handler.config.MiniIOConnection.BucketName, miniioFilename)
		if err != nil {
			logger.Log.Info(err.Error())
			http.Error(res, "Select binary data error", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		info, err := file.Stat()
		if err != nil {
			logger.Log.Info(err.Error())
			http.Error(res, "Error stating object", http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
		res.Header().Set("Content-Type", "application/octet-stream")
		res.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size))

		if _, err = io.Copy(res, file); err != nil {
			logger.Log.Info(err.Error())
			http.Error(res, "Failed to stream file", http.StatusInternalServerError)
		}
	}
}

// DeleteBinary удаляет бинарные данные из miniio.
func (handler *BinaryHandler) DeleteBinary() http.HandlerFunc {
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

		filename := string(body)
		miniioFilename := fmt.Sprintf("%s/%s", userID, filename)

		err = handler.service.DeleteBinary(req.Context(), handler.config.MiniIOConnection.BucketName, miniioFilename)
		if err != nil {
			logger.Log.Info(err.Error())
			http.Error(res, "Delete binary data error", http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusOK)
		io.WriteString(res, "Binary data "+filename+" deleted successfully")
	}
}
