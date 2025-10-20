package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/nu-kotov/GophKeeper/internal/config"
	"github.com/nu-kotov/GophKeeper/internal/dberrors"
	"github.com/nu-kotov/GophKeeper/internal/mocks"
	"github.com/nu-kotov/GophKeeper/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestCredentialsHandler_AddCredentials_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockCredentialsStorage(ctrl)
	mockStorage.EXPECT().InsertCredentialsData(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	r := chi.NewRouter()

	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	config.SecretKey = testSecret

	NewCredentialsHandler(r, config, mockStorage)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/credentials/add", map[string]string{
		"data_id":  testDataID,
		"login":    testLogin,
		"password": testEncriptedPass,
	},
		true,
	)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 201, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "Credentials "+testDataID+" added successfully", string(body), "Response text didn't match expected")
}

func TestCredentialsHandler_AddCredentials_AlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockCredentialsStorage(ctrl)
	mockStorage.EXPECT().InsertCredentialsData(gomock.Any(), gomock.Any(), gomock.Any()).Return(dberrors.ErrConflict)

	r := chi.NewRouter()

	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	config.SecretKey = testSecret

	NewCredentialsHandler(r, config, mockStorage)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/credentials/add", map[string]string{
		"data_id":  testDataID,
		"login":    testLogin,
		"password": testEncriptedPass,
	},
		true,
	)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 409, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "Credentials "+testDataID+" already exists", string(body), "Response text didn't match expected")
}

func TestCredentialsHandler_AddCredentials_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockCredentialsStorage(ctrl)

	r := chi.NewRouter()
	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	NewCredentialsHandler(r, config, mockStorage)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/credentials/add", map[string]string{
		"data_id":  testDataID,
		"login":    testLogin,
		"password": testEncriptedPass,
	},
		false)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 401, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "Please log in", string(body), "Response text didn't match expected")
}

func TestCredentialsHandler_GetCredentials_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCreds := models.Credentials{
		DataID:   testDataID,
		Login:    testLogin,
		Password: testEncriptedPass,
	}
	mockStorage := mocks.NewMockCredentialsStorage(ctrl)
	mockStorage.EXPECT().SelectCredentialsData(gomock.Any(), gomock.Any(), gomock.Any()).Return(&testCreds, nil)

	r := chi.NewRouter()

	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	config.SecretKey = testSecret

	NewCredentialsHandler(r, config, mockStorage)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/credentials/get", map[string]string{
		"data_id": testDataID,
	},
		true,
	)

	var respCreds models.Credentials
	err = json.NewDecoder(resp.Body).Decode(&respCreds)
	assert.NoError(t, err, "error response decoding")

	assert.Equal(t, 200, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, testCreds, respCreds, "Response didn't match expected")
}

func TestCredentialsHandler_GetCredentials_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockCredentialsStorage(ctrl)
	mockStorage.EXPECT().SelectCredentialsData(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, dberrors.ErrNotFound)

	r := chi.NewRouter()

	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	config.SecretKey = testSecret

	NewCredentialsHandler(r, config, mockStorage)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/credentials/get", map[string]string{
		"data_id": testDataID,
	},
		true,
	)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 404, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "Credentials "+testDataID+" not found", string(body), "Response text didn't match expected")
}

func TestCredentialsHandler_GetCredentials_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockCredentialsStorage(ctrl)

	r := chi.NewRouter()
	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	NewCredentialsHandler(r, config, mockStorage)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/credentials/get", map[string]string{
		"data_id": testDataID,
	},
		false)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 401, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "Please log in", string(body), "Response text didn't match expected")
}

func TestCredentialsHandler_DeleteCredentials_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockCredentialsStorage(ctrl)
	mockStorage.EXPECT().DeleteCredentialsData(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	r := chi.NewRouter()

	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	config.SecretKey = testSecret

	NewCredentialsHandler(r, config, mockStorage)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/credentials/delete", map[string]string{
		"data_id": testDataID,
	},
		true,
	)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "Credentials "+testDataID+" deleted successfully", string(body), "Response text didn't match expected")
}

func TestCredentialsHandler_DeleteCredentials_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockCredentialsStorage(ctrl)

	r := chi.NewRouter()
	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	NewCredentialsHandler(r, config, mockStorage)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/credentials/delete", map[string]string{
		"data_id": testDataID,
	},
		false)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 401, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "Please log in", string(body), "Response text didn't match expected")
}
