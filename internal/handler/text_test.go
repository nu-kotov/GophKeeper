package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/nu-kotov/GophKeeper/internal/auth"
	"github.com/nu-kotov/GophKeeper/internal/config"
	"github.com/nu-kotov/GophKeeper/internal/dberrors"
	"github.com/nu-kotov/GophKeeper/internal/mocks"
	"github.com/stretchr/testify/assert"
)

var (
	testUserID string = "974af87d-210e-41d3-b93c-2f938685862f"
	testSecret string = "testkey"
	testText   string = "Test text"
	testDataID string = "test_text"
)

func makeRequest(t *testing.T, ts *httptest.Server, method, path string, body any, isAuth bool) *http.Response {
	var reqBody *bytes.Reader

	if body != nil {
		b, err := json.Marshal(body)
		assert.NoError(t, err, "failed to marshal body")

		reqBody = bytes.NewReader(b)
	} else {
		reqBody = bytes.NewReader(nil)
	}

	req, err := http.NewRequest(method, ts.URL+path, reqBody)
	assert.NoError(t, err, "failed to create request")

	req.Header.Set("Content-Type", "application/json")
	if isAuth == true {
		token, _ := auth.BuildJWTString(
			testUserID,
			"test_login",
			time.Hour*72,
			testSecret,
		)
		cookie := &http.Cookie{
			Name:     "token",
			Value:    token,
			HttpOnly: true,
			Path:     "/",
		}
		req.AddCookie(cookie)
	}

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err, "request failed")

	return resp
}

func TestTextHandler_AddText_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockTextStorage(ctrl)
	mockStorage.EXPECT().InsertTextData(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	r := chi.NewRouter()

	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	config.SecretKey = testSecret

	NewTextHandler(r, config, mockStorage)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/text/add", map[string]string{
		"data_id": testDataID,
		"text":    testText,
	},
		true,
	)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 201, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "Text "+testDataID+" added successfully", string(body), "Response text didn't match expected")
}

func TestTextHandler_AddText_AlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockTextStorage(ctrl)
	mockStorage.EXPECT().InsertTextData(gomock.Any(), gomock.Any(), gomock.Any()).Return(dberrors.ErrConflict)

	r := chi.NewRouter()

	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	config.SecretKey = testSecret

	NewTextHandler(r, config, mockStorage)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/text/add", map[string]string{
		"data_id": testDataID,
		"text":    testText,
	},
		true,
	)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 409, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "Text "+testDataID+" already exists", string(body), "Response text didn't match expected")
}

func TestTextHandler_AddText_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockTextStorage(ctrl)

	r := chi.NewRouter()
	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	NewTextHandler(r, config, mockStorage)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/text/add", map[string]string{
		"data_id": testDataID,
		"text":    testText,
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
	mockStorage.EXPECT().SelectTextData(gomock.Any(), gomock.Any(), gomock.Any()).Return(testText, nil)

	r := chi.NewRouter()

	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	config.SecretKey = testSecret

	NewTextHandler(r, config, mockStorage)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/text/get", map[string]string{
		"data_id": testDataID,
	},
		true,
	)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, testText, string(body), "Response text didn't match expected")
}

func TestTextHandler_GetText_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockTextStorage(ctrl)
	mockStorage.EXPECT().SelectTextData(gomock.Any(), gomock.Any(), gomock.Any()).Return("", dberrors.ErrNotFound)

	r := chi.NewRouter()

	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	config.SecretKey = testSecret

	NewTextHandler(r, config, mockStorage)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/text/get", map[string]string{
		"data_id": testDataID,
	},
		true,
	)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 404, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "Text "+testDataID+" not found", string(body), "Response text didn't match expected")
}

func TestTextHandler_GetText_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockTextStorage(ctrl)

	r := chi.NewRouter()
	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	NewTextHandler(r, config, mockStorage)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/text/get", map[string]string{
		"data_id": testDataID,
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

	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	config.SecretKey = testSecret

	NewTextHandler(r, config, mockStorage)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/text/delete", map[string]string{
		"data_id": testDataID,
	},
		true,
	)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "Text "+testDataID+" deleted successfully", string(body), "Response text didn't match expected")
}

func TestTextHandler_DeleteText_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockTextStorage(ctrl)

	r := chi.NewRouter()
	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	NewTextHandler(r, config, mockStorage)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/text/delete", map[string]string{
		"data_id": testDataID,
	},
		false)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 401, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "Please log in", string(body), "Response text didn't match expected")
}
