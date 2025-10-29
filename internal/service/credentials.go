package service

import (
	"context"

	"github.com/nu-kotov/GophKeeper/internal/keeper_errors"
	"github.com/nu-kotov/GophKeeper/internal/models"
)

// CredentialsStorage — интерфейс, который уже реализован в твоем storage.
type CredentialsStorage interface {
	InsertCredentialsData(context.Context, string, *models.Credentials) error
	SelectCredentialsData(context.Context, string, string) (*models.Credentials, error)
	DeleteCredentialsData(context.Context, string, string) error
}

// CredentialsService — структура сервиса для работы с кредами.
type CredentialsService struct {
	storage CredentialsStorage
}

// NewCredentialsService — конструктор сервиса для работы с кредами.
func NewCredentialsService(storage CredentialsStorage) *CredentialsService {
	return &CredentialsService{storage: storage}
}

// AddCredentials — добавление кредов в бд.
func (s *CredentialsService) AddCredentials(ctx context.Context, userID string, creds *models.Credentials) error {
	if creds.DataID == "" || creds.Login == "" {
		return keeper_errors.ErrInvalidCredentials
	}
	return s.storage.InsertCredentialsData(ctx, userID, creds)
}

// GetCredentials — получение кредов из бд.
func (s *CredentialsService) GetCredentials(ctx context.Context, userID, dataID string) (*models.Credentials, error) {
	return s.storage.SelectCredentialsData(ctx, userID, dataID)
}

// DeleteCredentials — удаление кредов из бд.
func (s *CredentialsService) DeleteCredentials(ctx context.Context, userID, dataID string) error {
	return s.storage.DeleteCredentialsData(ctx, userID, dataID)
}
