package client_service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nu-kotov/GophKeeper/internal/models"
	"github.com/phedde/luhn-algorithm"
)

// CardService - структура сервиса клиента для работы с данными банковских карт.
type CardService struct {
	BaseURL    string
	HTTPClient *http.Client
	Key        []byte
	LoadCookie func() (*http.Cookie, error)
	Encrypt    func([]byte, string) (string, error)
	Decrypt    func([]byte, string) (string, error)
}

// CardPayload - структура запроса с данными банковской карты.
type CardPayload struct {
	Number string `json:"number"`
	Expiry string `json:"expiry"`
	CVV    string `json:"cvv"`
	Name   string `json:"name"`
}

// AddCard - сохраняет данные карты.
func (s *CardService) AddCard(id, number, expiry, cvv, holder string) (string, error) {

	if err := validateCardData(number, expiry, cvv, holder); err != nil {
		return "", err
	}

	card := models.CardPayload{Number: number, Expiry: expiry, CVV: cvv, Name: holder}
	cardJSON, err := json.Marshal(card)
	if err != nil {
		return "", err
	}

	encryptedCardData, err := s.Encrypt([]byte(s.Key), string(cardJSON))
	if err != nil {
		return "", err
	}

	payload := map[string]string{
		"data_id": id,
		"card":    encryptedCardData,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	cookie, err := s.LoadCookie()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", s.BaseURL+"/api/card/add", bytes.NewBuffer(body))
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

// GetCard - получает данные карты.
func (s *CardService) GetCard(id string) (*CardPayload, error) {
	payload := map[string]string{
		"data_id": id,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	cookie, err := s.LoadCookie()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", s.BaseURL+"/api/card/get", bytes.NewBuffer(body))
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

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	decryptedCardData, err := s.Decrypt([]byte(s.Key), string(respBody))
	if err != nil {
		return nil, err
	}

	var card CardPayload
	err = json.Unmarshal([]byte(decryptedCardData), &card)
	if err != nil {
		return nil, err
	}

	return &card, nil
}

// GetCard - удаляет данные карты.
func (s *CardService) DelCard(id string) (string, error) {
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

	req, err := http.NewRequest("POST", s.BaseURL+"/api/card/delete", bytes.NewBuffer(body))
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

func validateCardData(number, expiry, cvv, holder string) error {
	// Проверка номера карты (Luhn + тип карты)
	intNumber, err := strconv.ParseInt(number, 10, 64)
	if err != nil {
		return err
	}
	if isValid := luhn.IsValid(intNumber); !isValid {
		return errors.New("invalid card number")
	}

	// CVV: ровно 3 цифры
	matched, _ := regexp.MatchString(`^\d{3}$`, cvv)
	if !matched {
		return errors.New("CVV must be exactly 3 digits")
	}

	// Expiry: не просрочен
	month, year, err := parseExpiryDate(expiry)
	if err != nil {
		return err
	}
	if isExpired(month, year) {
		return errors.New("card is expired")
	}

	if strings.TrimSpace(holder) == "" {
		return errors.New("holder is required")
	}

	return nil
}

func parseExpiryDate(expiry string) (int, int, error) {
	expiry = strings.TrimSpace(expiry)
	parts := strings.Split(expiry, "/")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("expiry must be in format MM/YY or MM/YYYY")
	}

	month, err := strconv.Atoi(parts[0])
	if err != nil || month < 1 || month > 12 {
		return 0, 0, fmt.Errorf("invalid month")
	}

	year, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid year")
	}

	if year < 100 {
		year += 2000
	}

	return month, year, nil
}

func isExpired(month, year int) bool {
	now := time.Now()
	exp := time.Date(year, time.Month(month)+1, 0, 23, 59, 59, 0, time.UTC)
	return now.After(exp)
}
