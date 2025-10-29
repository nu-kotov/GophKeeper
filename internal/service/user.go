package service

import (
	"context"

	"github.com/nu-kotov/GophKeeper/internal/models"
)

// UsersStorage - интерфейс хранилища для работы с пользователями.
type UsersStorage interface {
	InsertUserData(context.Context, *models.UserData) error
	SelectUserData(context.Context, *models.UserData) (*models.UserData, error)
}

// UsersService — структура сервиса для работы с пользователями.
type UsersService struct {
	storage UsersStorage
}

// NewUsersService — конструктор сервиса для работы с пользователями.
func NewUsersService(storage UsersStorage) *UsersService {
	return &UsersService{storage: storage}
}

// AddUser — добавление пользователя в бд.
func (s *UsersService) AddUser(ctx context.Context, user *models.UserData) error {
	return s.storage.InsertUserData(ctx, user)
}

// GetUser — получение пользователя из бд.
func (s *UsersService) GetUser(ctx context.Context, user *models.UserData) (*models.UserData, error) {
	return s.storage.SelectUserData(ctx, user)
}
