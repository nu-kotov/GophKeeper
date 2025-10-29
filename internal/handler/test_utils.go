package handler

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nu-kotov/GophKeeper/internal/auth"
	"github.com/nu-kotov/GophKeeper/internal/testvars"
	"github.com/stretchr/testify/assert"
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
	if isAuth {
		add_token_cookie(req)
	}

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err, "request failed")

	return resp
}

func uploadBinaryFile(t *testing.T, ts *httptest.Server, method, path, filename string, data []byte, isAuth bool) *http.Response {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", filename)
	assert.NoError(t, err, "failed to create form file")
	part.Write(data)
	writer.Close()

	req, err := http.NewRequest(method, ts.URL+path, &body)
	assert.NoError(t, err, "failed to create request")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-Filename", filename)
	if isAuth {
		add_token_cookie(req)
	}

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err, "upload request failed")

	return resp
}

func add_token_cookie(req *http.Request) {
	token, _ := auth.BuildJWTString(
		testvars.TestUserID,
		"test_login",
		time.Hour*72,
		testvars.TestSecret,
	)
	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
	}
	req.AddCookie(cookie)
}
