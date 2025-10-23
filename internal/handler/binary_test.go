package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/minio/minio-go/v7"
	"github.com/nu-kotov/GophKeeper/internal/config"
	"github.com/nu-kotov/GophKeeper/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestBinaryHandler_AddBinary_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUploadInfo := minio.UploadInfo{Bucket: testUserID}

	mockStorage := mocks.NewMockBinaryStorage(ctrl)
	mockStorage.EXPECT().InsertBinaryData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&testUploadInfo, nil)

	r := chi.NewRouter()

	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	config.SecretKey = testSecret

	NewBinaryHandler(r, config, mockStorage)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := uploadBinaryFile(t, ts, http.MethodPost, "/api/binary/add", testFileName, testBinaryData, true)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 201, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "File "+testFileName+" uploaded successfully", string(body), "Response text didn't match expected")
}

func TestBinaryHandler_AddBinary_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockBinaryStorage(ctrl)

	r := chi.NewRouter()
	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	NewBinaryHandler(r, config, mockStorage)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := uploadBinaryFile(t, ts, http.MethodPost, "/api/binary/add", testFileName, testBinaryData, false)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 401, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "Please log in", string(body), "Response text didn't match expected")
}

func TestBinaryHandler_GetBinary_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockBinaryStorage(ctrl)

	r := chi.NewRouter()
	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	NewBinaryHandler(r, config, mockStorage)

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest(http.MethodPost, ts.URL+"/api/binary/get", strings.NewReader(testDataID))
	assert.NoError(t, err, "failed to create request")

	req.Header.Add("Content-Type", "text/plain")

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err, "request failed")
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, 401, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "Please log in", string(body), "Response text didn't match expected")
}

func TestBinaryHandler_DeleteBinary_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockBinaryStorage(ctrl)
	mockStorage.EXPECT().DeleteBinaryData(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	r := chi.NewRouter()

	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	config.SecretKey = testSecret

	NewBinaryHandler(r, config, mockStorage)

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest(http.MethodPost, ts.URL+"/api/binary/delete", strings.NewReader(testDataID))
	assert.NoError(t, err, "failed to create request")

	add_token_cookie(req)
	req.Header.Add("Content-Type", "text/plain")

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err, "request failed")
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "Binary data "+testDataID+" deleted successfully", string(body), "Response text didn't match expected")
}

func TestBinaryHandler_DeleteBinary_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockBinaryStorage(ctrl)

	r := chi.NewRouter()
	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	NewBinaryHandler(r, config, mockStorage)

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest(http.MethodPost, ts.URL+"/api/binary/delete", strings.NewReader(testDataID))
	assert.NoError(t, err, "failed to create request")

	req.Header.Add("Content-Type", "text/plain")

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err, "request failed")
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, 401, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "Please log in", string(body), "Response text didn't match expected")
}
