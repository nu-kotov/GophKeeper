package handler

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/nu-kotov/GophKeeper/internal/config"
	"github.com/nu-kotov/GophKeeper/internal/middleware"
	"github.com/nu-kotov/GophKeeper/internal/models"
)

type TextStorage interface {
	InsertTextData(context.Context, *models.TextData) error
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
	return func(res http.ResponseWriter, req *http.Request) {}
}

func (handler *TextHandler) GetText() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {}
}

func (handler *TextHandler) DeleteText() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {}
}
