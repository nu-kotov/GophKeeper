package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/nu-kotov/GophKeeper/internal/client"
	"github.com/nu-kotov/GophKeeper/internal/config"
	"github.com/nu-kotov/GophKeeper/internal/keeper_errors"
	"github.com/nu-kotov/GophKeeper/internal/mocks"
	"github.com/nu-kotov/GophKeeper/internal/models"
	"github.com/nu-kotov/GophKeeper/internal/service"
	"github.com/nu-kotov/GophKeeper/internal/testvars"
	"github.com/stretchr/testify/assert"
)

func TestCardHandler_AddCard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockCardStorage(ctrl)
	mockStorage.EXPECT().InsertCardData(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	r := chi.NewRouter()

	service := service.NewCardService(mockStorage)

	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	config.SecretKey = testvars.TestSecret

	NewCardHandler(r, config, service)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/card/add", map[string]string{
		"data_id": testvars.TestDataID,
		"card":    testvars.TestEncriptedCard,
	},
		true,
	)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 201, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "Card "+testvars.TestDataID+" added successfully", string(body), "Response text didn't match expected")
}

func TestCardHandler_AddCard_AlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockCardStorage(ctrl)
	mockStorage.EXPECT().InsertCardData(gomock.Any(), gomock.Any(), gomock.Any()).Return(keeper_errors.ErrConflict)

	r := chi.NewRouter()

	service := service.NewCardService(mockStorage)

	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	config.SecretKey = testvars.TestSecret

	NewCardHandler(r, config, service)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/card/add", map[string]string{
		"data_id": testvars.TestDataID,
		"card":    testvars.TestEncriptedCard,
	},
		true,
	)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 409, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "Card "+testvars.TestDataID+" already exists", string(body), "Response text didn't match expected")
}

func TestCardHandler_AddCard_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockCardStorage(ctrl)

	r := chi.NewRouter()

	service := service.NewCardService(mockStorage)

	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	NewCardHandler(r, config, service)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/card/add", map[string]string{
		"data_id": testvars.TestDataID,
		"card":    testvars.TestEncriptedCard,
	},
		false)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 401, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "Please log in", string(body), "Response text didn't match expected")
}

func TestCardHandler_GetCard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCard := models.CardPayload{
		Number: testvars.TestCardNumber,
		Expiry: testvars.TestCardExp,
		CVV:    testvars.TestCVV,
		Name:   testvars.TestUser,
	}
	cardJSON, err := json.Marshal(testCard)
	assert.NoError(t, err, "error test card marshalling")

	testEncryptedCardData, err := client.Encrypt([]byte(testvars.TestClientKey), string(cardJSON))
	assert.NoError(t, err, "error test card encripting")

	mockStorage := mocks.NewMockCardStorage(ctrl)
	mockStorage.EXPECT().SelectCardData(gomock.Any(), gomock.Any(), gomock.Any()).Return(testEncryptedCardData, nil)

	r := chi.NewRouter()

	service := service.NewCardService(mockStorage)

	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	config.SecretKey = testvars.TestSecret

	NewCardHandler(r, config, service)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/card/get", map[string]string{
		"data_id": testvars.TestDataID,
	},
		true,
	)

	respBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "error response reading")

	decryptedCardData, err := client.Decrypt([]byte(testvars.TestClientKey), string(respBody))
	assert.NoError(t, err, "error response decripting")

	var respCardData models.CardPayload
	err = json.Unmarshal([]byte(decryptedCardData), &respCardData)
	assert.NoError(t, err, "error response unmarshalling")

	assert.Equal(t, 200, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, testCard, respCardData, "Response didn't match expected")
}

func TestCardHandler_GetCard_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockCardStorage(ctrl)
	mockStorage.EXPECT().SelectCardData(gomock.Any(), gomock.Any(), gomock.Any()).Return("", keeper_errors.ErrNotFound)

	r := chi.NewRouter()

	service := service.NewCardService(mockStorage)

	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	config.SecretKey = testvars.TestSecret

	NewCardHandler(r, config, service)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/card/get", map[string]string{
		"data_id": testvars.TestDataID,
	},
		true,
	)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 404, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "Card "+testvars.TestDataID+" not found", string(body), "Response text didn't match expected")
}

func TestCardHandler_GetCard_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockCardStorage(ctrl)

	r := chi.NewRouter()

	service := service.NewCardService(mockStorage)

	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	NewCardHandler(r, config, service)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/card/get", map[string]string{
		"data_id": testvars.TestDataID,
	},
		false)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 401, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "Please log in", string(body), "Response text didn't match expected")
}

func TestCardHandler_DeleteCard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockCardStorage(ctrl)
	mockStorage.EXPECT().DeleteCardData(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	r := chi.NewRouter()

	service := service.NewCardService(mockStorage)

	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	config.SecretKey = testvars.TestSecret

	NewCardHandler(r, config, service)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/card/delete", map[string]string{
		"data_id": testvars.TestDataID,
	},
		true,
	)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "Card "+testvars.TestDataID+" deleted successfully", string(body), "Response text didn't match expected")
}

func TestCardHandler_DeleteCard_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockCardStorage(ctrl)

	r := chi.NewRouter()

	service := service.NewCardService(mockStorage)

	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	NewCardHandler(r, config, service)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/card/delete", map[string]string{
		"data_id": testvars.TestDataID,
	},
		false)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 401, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "Please log in", string(body), "Response text didn't match expected")
}
