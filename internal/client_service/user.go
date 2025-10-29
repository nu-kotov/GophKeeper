package client_service

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/nu-kotov/GophKeeper/internal/keeper_errors"
)

// UserService - структура сервиса клиента для работы с пользователями.
type UserService struct {
	BaseURL    string
	HTTPClient *http.Client
	Key        []byte
	LoadCookie func() (*http.Cookie, error)
	SaveCookie func(*http.Cookie) error
	Encrypt    func([]byte, string) (string, error)
	Decrypt    func([]byte, string) (string, error)
}

// Login - авторизует пользователя.
func (s *UserService) Login(login, password string) (string, error) {

	loginData := map[string]string{
		"login":    login,
		"password": password,
	}

	body, err := json.Marshal(loginData)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(s.BaseURL+"/api/user/login", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "token" {
			err := s.SaveCookie(cookie)
			if err != nil {
				return "", err
			}
			return "Login successful. Session saved.", nil
		}
	}

	return "", keeper_errors.ErrTokenNotFound
}

// Register - регистрирует нового пользователя.
func (s *UserService) Register(login, password string) (string, error) {

	loginData := map[string]string{
		"login":    login,
		"password": password,
	}

	body, err := json.Marshal(loginData)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(s.BaseURL+"/api/user/register", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "token" {
			err := s.SaveCookie(cookie)
			if err != nil {
				return "", err
			}
			return "Register successful. Session saved.", nil
		}
	}

	return "", keeper_errors.ErrTokenNotFound
}
