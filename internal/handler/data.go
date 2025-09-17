package handler

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/nu-kotov/GophKeeper/internal/config"
	"github.com/nu-kotov/GophKeeper/internal/middleware"
	"github.com/nu-kotov/GophKeeper/internal/models"
)

type DataStorage interface {
	InsertData(context.Context, *models.PrivateData) error
	SelectDataByID(context.Context, string) (*models.PrivateData, error)
	SelectAllDataByUserID(context.Context) ([]models.PrivateData, error)
	DeleteDataByID(context.Context, string) error
}

type DataHandler struct {
	Config  *config.Config
	Storage DataStorage
}

func NewDataHandler(router *chi.Mux, cfg *config.Config, storage DataStorage) {

	handler := &DataHandler{
		Config:  cfg,
		Storage: storage,
	}

	middlewareStack := middleware.Chain(
		middleware.RequestLogger,
	)

	router.Post(`/api/data/add`, middlewareStack(handler.AddData()))
	router.Post(`/api/data/get`, middlewareStack(handler.GetDataByID()))
	router.Post(`/api/data/get-all`, middlewareStack(handler.GetAllDataByUserID()))
	router.Delete(`/api/data/delete`, middlewareStack(handler.DeleteDataByID()))

}

func (handler *DataHandler) AddData() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {}
}

func (handler *DataHandler) GetDataByID() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {}
}

func (handler *DataHandler) GetAllDataByUserID() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {}
}

func (handler *DataHandler) DeleteDataByID() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {}
}
