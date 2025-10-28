package service

import (
	"context"

	"github.com/nu-kotov/GophKeeper/internal/models"
)

// TextStorage - интерфейс хранилища для работы с текстом.
type TextStorage interface {
	InsertTextData(context.Context, string, *models.TextData) error
	SelectTextData(context.Context, string, string) (string, error)
	DeleteTextData(context.Context, string, string) error
}

// TextService — структура сервиса для работы с текстом.
type TextService struct {
	storage TextStorage
}

// NewTextService — конструктор сервиса для работы с текстом.
func NewTextService(storage TextStorage) *TextService {
	return &TextService{storage: storage}
}

// AddText — добавление текста в бд.
func (s *TextService) AddText(ctx context.Context, userID string, txt *models.TextData) error {
	return s.storage.InsertTextData(ctx, userID, txt)
}

// GetText — получение текста из бд.
func (s *TextService) GetText(ctx context.Context, userID, dataID string) (string, error) {
	return s.storage.SelectTextData(ctx, userID, dataID)
}

// DeleteText — удаление текста из бд.
func (s *TextService) DeleteText(ctx context.Context, userID, dataID string) error {
	return s.storage.DeleteTextData(ctx, userID, dataID)
}
