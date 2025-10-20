package handler

import (
	"bytes"
	"encoding/json"
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
	testText          string = "Test text"
	testDataID        string = "test_text"
	testLogin         string = "test_login"
	testEncriptedPass string = "$argon2id$v=19$m=65536,t=1,p=20$usM+ExLZFSqBWTWUffEd9w$Sm+1/yidgUWG1AnvGeCbgEFvo3MKCBf0uig2Wcuty60"
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
