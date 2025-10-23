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
	"github.com/stretchr/testify/assert"
)

var (
	testUserID        string = "974af87d-210e-41d3-b93c-2f938685862f"
	testSecret        string = "testkey"
	testClientKey     string = "12345678901234567890123456789012"
	testCardNumber    string = "4433062851071851"
	testCardExp       string = "2/2027"
	testCVV           string = "333"
	testText          string = "Test text"
	testBinaryData    []byte = []byte("test binary")
	testDataID        string = "test_data"
	testUser          string = "test_login"
	testFileName      string = "test_file.txt"
	testEncriptedPass string = "$argon2id$v=19$m=65536,t=1,p=20$usM+ExLZFSqBWTWUffEd9w$Sm+1/yidgUWG1AnvGeCbgEFvo3MKCBf0uig2Wcuty60"
	testEncriptedCard string = "37s9XmpC4zJZsajcuub7burtXja9pD2qRBNBz8F/MYOYtEdGrHchdp6eRcX2d9c7H24HhVgMFVcjz29Hye3tW0VSqUs0g8QWAF2wDlBmCgmHCsuVpZT1cis60ZRZMqR7DvGBWSmppS2vzH4="
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
