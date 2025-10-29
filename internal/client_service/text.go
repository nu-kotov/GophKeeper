package client_service

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// TextService - структура сервиса клиента для работы с текстовыми данными.
type TextService struct {
	BaseURL    string
	HTTPClient *http.Client
	Key        []byte
	LoadCookie func() (*http.Cookie, error)
	Encrypt    func([]byte, string) (string, error)
	Decrypt    func([]byte, string) (string, error)
}

// TextData - структура запроса с текстом.
type TextData struct {
	DataID string `json:"data_id"`
	Text   string `json:"text"`
}

// AddText - сохраняет текст.
func (s *TextService) AddText(id, txt string) (string, error) {

	data := TextData{
		DataID: id,
		Text:   txt,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	cookie, err := s.LoadCookie()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", s.BaseURL+"/api/text/add", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.AddCookie(cookie)
	req.Header.Add("Content-Type", "application/json")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

// GetCard - получает данные карты.
func (s *TextService) GetText(id string) (string, error) {

	payload := map[string]string{
		"data_id": id,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	cookie, err := s.LoadCookie()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", s.BaseURL+"/api/text/get", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.AddCookie(cookie)
	req.Header.Add("Content-Type", "application/json")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

// DelText - удаляет текст.
func (s *TextService) DelText(id string) (string, error) {

	payload := map[string]string{
		"data_id": id,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	cookie, err := s.LoadCookie()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", s.BaseURL+"/api/text/delete", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.AddCookie(cookie)
	req.Header.Add("Content-Type", "application/json")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}
