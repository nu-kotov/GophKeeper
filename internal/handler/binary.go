package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/minio/minio-go/v7"

	"github.com/nu-kotov/GophKeeper/internal/auth"
	"github.com/nu-kotov/GophKeeper/internal/config"
	"github.com/nu-kotov/GophKeeper/internal/logger"
	"github.com/nu-kotov/GophKeeper/internal/middleware"
)

type BinaryStorage interface {
	InsertBinaryData(context.Context, string, string, io.ReadCloser) (*minio.UploadInfo, error)
	SelectBinaryData(context.Context, string, string) (*minio.Object, error)
	// DeleteBinaryData(context.Context, string, string) error
}

type BinaryHandler struct {
	Config  *config.Config
	Storage BinaryStorage
}

func NewBinaryHandler(router *chi.Mux, cfg *config.Config, storage BinaryStorage) {

	handler := &BinaryHandler{
		Config:  cfg,
		Storage: storage,
	}

	middlewareStack := middleware.Chain(
		middleware.RequestLogger,
	)

	router.Post(`/api/binary/add`, middlewareStack(handler.AddBinary()))
	router.Post(`/api/binary/get`, middlewareStack(handler.GetBinary()))
	// router.Post(`/api/binary/delete`, middlewareStack(handler.DeleteBinary()))

}

func (handler *BinaryHandler) AddBinary() http.HandlerFunc {
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

		filename := req.Header.Get("X-Filename")
		if filename == "" {
			http.Error(res, "Missing X-Filename header", http.StatusBadRequest)
			return
		}

		miniioFilename := fmt.Sprintf("%s/%s", userID, filename)

		_, err = handler.Storage.InsertBinaryData(
			req.Context(),
			handler.Config.MiniIOConnection.BucketName,
			miniioFilename,
			req.Body,
		)
		if err != nil {
			logger.Log.Info(err.Error())
			http.Error(res, "File uploaded error", http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusOK)
		io.WriteString(res, "File uploaded")
	}
}

func (handler *BinaryHandler) GetBinary() http.HandlerFunc {
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

		filename := string(body)
		miniioFilename := fmt.Sprintf("%s/%s", userID, filename)

		file, err := handler.Storage.SelectBinaryData(req.Context(), handler.Config.MiniIOConnection.BucketName, miniioFilename)
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

// func (handler *BinaryHandler) DeleteBinary() http.HandlerFunc {
// 	return func(res http.ResponseWriter, req *http.Request) {
// 		token, err := req.Cookie("token")

// 		if err != nil {
// 			logger.Log.Info(err.Error())
// 			res.WriteHeader(http.StatusUnauthorized)
// 			return
// 		}

// 		userID, err := auth.GetUserID(token.Value, handler.Config.SecretKey)
// 		if err != nil {
// 			logger.Log.Info(err.Error())
// 			http.Error(res, err.Error(), http.StatusBadRequest)
// 			return
// 		}

// 		body, err := io.ReadAll(req.Body)
// 		if err != nil {
// 			logger.Log.Info(err.Error())
// 			http.Error(res, "Invalid body", http.StatusBadRequest)
// 			return
// 		}

// 		var jsonBody models.BinaryID
// 		if err = json.Unmarshal(body, &jsonBody); err != nil {
// 			http.Error(res, err.Error(), http.StatusBadRequest)
// 			return
// 		}

// 		err = handler.Storage.DeleteBinaryData(req.Context(), userID, jsonBody.DataID)
// 		if err != nil {
// 			logger.Log.Info(err.Error())
// 			http.Error(res, "Delete text data error", http.StatusInternalServerError)
// 			return
// 		}

// 		res.Header().Set("Content-Type", "text/plain")
// 		res.WriteHeader(http.StatusOK)
// 	}
// }
