package service

import (
	"context"

	"github.com/nu-kotov/GophKeeper/internal/models"
)

// CardStorage - интерфейс хранилища для работы с данными банковских карт.
type CardStorage interface {
	InsertCardData(context.Context, string, *models.CardData) error
	SelectCardData(context.Context, string, string) (string, error)
	DeleteCardData(context.Context, string, string) error
}

// CardService — структура сервиса для работы с данными банковских карт.
type CardService struct {
	storage CardStorage
}

// NewCardService — конструктор сервиса для работы с данными банковских карт.
func NewCardService(storage CardStorage) *CardService {
	return &CardService{storage: storage}
}

// AddCredentials — добавление кредов в бд.
func (s *CardService) AddCard(ctx context.Context, userID string, creds *models.CardData) error {
	return s.storage.InsertCardData(ctx, userID, creds)
}

// GetCredentials — получение кредов из бд.
func (s *CardService) GetCard(ctx context.Context, userID, dataID string) (string, error) {
	return s.storage.SelectCardData(ctx, userID, dataID)
}

// DeleteCredentials — удаление кредов из бд.
func (s *CardService) DeleteCard(ctx context.Context, userID, dataID string) error {
	return s.storage.DeleteCardData(ctx, userID, dataID)
}
