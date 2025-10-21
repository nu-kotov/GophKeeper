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
	"github.com/nu-kotov/GophKeeper/internal/dberrors"
	"github.com/nu-kotov/GophKeeper/internal/mocks"
	"github.com/nu-kotov/GophKeeper/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestCardHandler_AddCard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockCardStorage(ctrl)
	mockStorage.EXPECT().InsertCardData(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	r := chi.NewRouter()

	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	config.SecretKey = testSecret

	NewCardHandler(r, config, mockStorage)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/card/add", map[string]string{
		"data_id": testDataID,
		"card":    testEncriptedCard,
	},
		true,
	)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 201, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "Card "+testDataID+" added successfully", string(body), "Response text didn't match expected")
}

func TestCardHandler_AddCard_AlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockCardStorage(ctrl)
	mockStorage.EXPECT().InsertCardData(gomock.Any(), gomock.Any(), gomock.Any()).Return(dberrors.ErrConflict)

	r := chi.NewRouter()

	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	config.SecretKey = testSecret

	NewCardHandler(r, config, mockStorage)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/card/add", map[string]string{
		"data_id": testDataID,
		"card":    testEncriptedCard,
	},
		true,
	)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 409, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "Card "+testDataID+" already exists", string(body), "Response text didn't match expected")
}

func TestCardHandler_AddCard_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockCardStorage(ctrl)

	r := chi.NewRouter()
	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	NewCardHandler(r, config, mockStorage)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/card/add", map[string]string{
		"data_id": testDataID,
		"card":    testEncriptedCard,
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
		Number: testCardNumber,
		Expiry: testCardExp,
		CVV:    testCVV,
		Name:   testUser,
	}
	cardJSON, err := json.Marshal(testCard)
	assert.NoError(t, err, "error test card marshalling")

	testEncryptedCardData, err := client.Encrypt([]byte(testClientKey), string(cardJSON))
	assert.NoError(t, err, "error test card encripting")

	mockStorage := mocks.NewMockCardStorage(ctrl)
	mockStorage.EXPECT().SelectCardData(gomock.Any(), gomock.Any(), gomock.Any()).Return(testEncryptedCardData, nil)

	r := chi.NewRouter()

	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	config.SecretKey = testSecret

	NewCardHandler(r, config, mockStorage)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/card/get", map[string]string{
		"data_id": testDataID,
	},
		true,
	)

	respBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "error response reading")

	decryptedCardData, err := client.Decrypt([]byte(testClientKey), string(respBody))
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
	mockStorage.EXPECT().SelectCardData(gomock.Any(), gomock.Any(), gomock.Any()).Return("", dberrors.ErrNotFound)

	r := chi.NewRouter()

	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	config.SecretKey = testSecret

	NewCardHandler(r, config, mockStorage)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/card/get", map[string]string{
		"data_id": testDataID,
	},
		true,
	)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 404, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "Card "+testDataID+" not found", string(body), "Response text didn't match expected")
}

func TestCardHandler_GetCard_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockCardStorage(ctrl)

	r := chi.NewRouter()
	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	NewCardHandler(r, config, mockStorage)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/card/get", map[string]string{
		"data_id": testDataID,
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

	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	config.SecretKey = testSecret

	NewCardHandler(r, config, mockStorage)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/card/delete", map[string]string{
		"data_id": testDataID,
	},
		true,
	)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "Card "+testDataID+" deleted successfully", string(body), "Response text didn't match expected")
}

func TestCardHandler_DeleteCard_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockCardStorage(ctrl)

	r := chi.NewRouter()
	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	NewCardHandler(r, config, mockStorage)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/card/delete", map[string]string{
		"data_id": testDataID,
	},
		false)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 401, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "Please log in", string(body), "Response text didn't match expected")
}
