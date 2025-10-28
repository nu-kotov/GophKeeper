package client_service

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type CredentialsService struct {
	BaseURL    string
	HTTPClient *http.Client
	Key        []byte
	LoadCookie func() (*http.Cookie, error)
	Encrypt    func([]byte, string) (string, error)
	Decrypt    func([]byte, string) (string, error)
}

type Credential struct {
	ID       string `json:"data_id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (s *CredentialsService) AddCreds(id, login, password string) (string, error) {

	encryptedPassword, err := s.Encrypt([]byte(s.Key), password)
	if err != nil {
		return "", err
	}

	payload := map[string]string{
		"data_id":  id,
		"login":    login,
		"password": encryptedPassword,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	cookie, err := s.LoadCookie()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", s.BaseURL+"/api/credentials/add", bytes.NewBuffer(body))
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

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}

func (s *CredentialsService) GetCreds(id string) (*Credential, error) {
	cookie, err := s.LoadCookie()
	if err != nil {
		return nil, err
	}

	payload := map[string]string{
		"data_id": id,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", s.BaseURL+"/api/credentials/get", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.AddCookie(cookie)
	req.Header.Add("Content-Type", "application/json")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result Credential

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	decryptedPassword, err := s.Decrypt(s.Key, result.Password)
	if err != nil {
		return nil, err
	}
	result.Password = decryptedPassword

	return &result, nil
}

func (s *CredentialsService) DelCreds(id string) (string, error) {
	cookie, err := s.LoadCookie()
	if err != nil {
		return "", err
	}

	payload := map[string]string{
		"data_id": id,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", s.BaseURL+"/api/credentials/delete", bytes.NewBuffer(body))
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

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}
