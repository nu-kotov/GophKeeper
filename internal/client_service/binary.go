package client_service

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// BinaryService - структура сервиса клиента для работы с бинарными данными.
type BinaryService struct {
	BaseURL    string
	HTTPClient *http.Client
	Key        []byte
	LoadCookie func() (*http.Cookie, error)
	Encrypt    func([]byte, string) (string, error)
	Decrypt    func([]byte, string) (string, error)
}

// AddBinary - сохраняет бинарные данные.
func (s *BinaryService) AddBinary(filePath string) (string, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	cookie, err := s.LoadCookie()
	if err != nil {
		return "", err
	}

	filename := filepath.Base(filePath)

	req, err := http.NewRequest("POST", s.BaseURL+"/api/binary/add", file)
	if err != nil {
		return "", err
	}

	req.AddCookie(cookie)
	req.Header.Add("X-Filename", filename)
	req.Header.Add("Content-Type", "application/octet-stream")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}

// GetBinary - получает бинарные данные.
func (s *BinaryService) GetBinary(id, filePath string) error {

	cookie, err := s.LoadCookie()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.BaseURL+"/api/binary/get", strings.NewReader(id))
	if err != nil {
		return err
	}

	req.AddCookie(cookie)
	req.Header.Add("Content-Type", "text/plain")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if filePath == "" {
		filePath = filepath.Base(id)
	}

	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// DelBinary - удаляет бинарные данные.
func (s *BinaryService) DelBinary(id string) (string, error) {
	cookie, err := s.LoadCookie()
	if err != nil {
		return "", nil
	}

	req, err := http.NewRequest("POST", s.BaseURL+"/api/binary/delete", strings.NewReader(id))
	if err != nil {
		return "", nil
	}

	req.AddCookie(cookie)
	req.Header.Add("Content-Type", "text/plain")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}
