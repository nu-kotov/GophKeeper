package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/nu-kotov/GophKeeper/internal/config"
	"github.com/nu-kotov/GophKeeper/internal/keeper_errors"
	"github.com/nu-kotov/GophKeeper/internal/mocks"
	"github.com/nu-kotov/GophKeeper/internal/service"
	"github.com/nu-kotov/GophKeeper/internal/testvars"
	"github.com/stretchr/testify/assert"
)

func TestTextHandler_AddText_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockTextStorage(ctrl)
	mockStorage.EXPECT().InsertTextData(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	r := chi.NewRouter()

	service := service.NewTextService(mockStorage)

	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	config.SecretKey = testvars.TestSecret

	NewTextHandler(r, config, service)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/text/add", map[string]string{
		"data_id": testvars.TestDataID,
		"text":    testvars.TestText,
	},
		true,
	)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 201, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "Text "+testvars.TestDataID+" added successfully", string(body), "Response text didn't match expected")
}

func TestTextHandler_AddText_AlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockTextStorage(ctrl)
	mockStorage.EXPECT().InsertTextData(gomock.Any(), gomock.Any(), gomock.Any()).Return(keeper_errors.ErrConflict)

	r := chi.NewRouter()

	service := service.NewTextService(mockStorage)

	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	config.SecretKey = testvars.TestSecret

	NewTextHandler(r, config, service)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/text/add", map[string]string{
		"data_id": testvars.TestDataID,
		"text":    testvars.TestText,
	},
		true,
	)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 409, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "Text "+testvars.TestDataID+" already exists", string(body), "Response text didn't match expected")
}

func TestTextHandler_AddText_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockTextStorage(ctrl)

	r := chi.NewRouter()

	service := service.NewTextService(mockStorage)

	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	NewTextHandler(r, config, service)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/text/add", map[string]string{
		"data_id": testvars.TestDataID,
		"text":    testvars.TestText,
	},
		false)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 401, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "Please log in", string(body), "Response text didn't match expected")
}

func TestTextHandler_GetText_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockTextStorage(ctrl)
	mockStorage.EXPECT().SelectTextData(gomock.Any(), gomock.Any(), gomock.Any()).Return(testvars.TestText, nil)

	r := chi.NewRouter()

	service := service.NewTextService(mockStorage)

	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	config.SecretKey = testvars.TestSecret

	NewTextHandler(r, config, service)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/text/get", map[string]string{
		"data_id": testvars.TestDataID,
	},
		true,
	)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, testvars.TestText, string(body), "Response text didn't match expected")
}

func TestTextHandler_GetText_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockTextStorage(ctrl)
	mockStorage.EXPECT().SelectTextData(gomock.Any(), gomock.Any(), gomock.Any()).Return("", keeper_errors.ErrNotFound)

	r := chi.NewRouter()

	service := service.NewTextService(mockStorage)

	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	config.SecretKey = testvars.TestSecret

	NewTextHandler(r, config, service)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/text/get", map[string]string{
		"data_id": testvars.TestDataID,
	},
		true,
	)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 404, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "Text "+testvars.TestDataID+" not found", string(body), "Response text didn't match expected")
}

func TestTextHandler_GetText_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockTextStorage(ctrl)

	r := chi.NewRouter()

	service := service.NewTextService(mockStorage)

	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	NewTextHandler(r, config, service)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/text/get", map[string]string{
		"data_id": testvars.TestDataID,
	},
		false)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 401, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "Please log in", string(body), "Response text didn't match expected")
}

func TestTextHandler_DeleteText_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockTextStorage(ctrl)
	mockStorage.EXPECT().DeleteTextData(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	r := chi.NewRouter()

	service := service.NewTextService(mockStorage)

	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	config.SecretKey = testvars.TestSecret

	NewTextHandler(r, config, service)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/text/delete", map[string]string{
		"data_id": testvars.TestDataID,
	},
		true,
	)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "Text "+testvars.TestDataID+" deleted successfully", string(body), "Response text didn't match expected")
}

func TestTextHandler_DeleteText_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockTextStorage(ctrl)

	r := chi.NewRouter()

	service := service.NewTextService(mockStorage)

	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	NewTextHandler(r, config, service)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/text/delete", map[string]string{
		"data_id": testvars.TestDataID,
	},
		false)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 401, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "Please log in", string(body), "Response text didn't match expected")
}
