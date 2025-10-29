package client

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nu-kotov/GophKeeper/internal/models"
	"github.com/nu-kotov/GophKeeper/internal/testvars"
	"github.com/stretchr/testify/assert"
)

func TestAddtxtCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err, "error read request body")

		var jsonBody models.TextData
		err = json.Unmarshal(body, &jsonBody)
		assert.NoError(t, err, "error body unmarshalling")

		assert.Equal(t, jsonBody.DataID, testvars.TestDataID, "Request data id didn't match expected")
		assert.Equal(t, jsonBody.Text, testvars.TestText, "Request text didn't match expected")

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, "Text "+jsonBody.DataID+" added successfully")
	}))
	defer server.Close()

	baseURL = server.URL

	svc := newTextService()

	resp, err := svc.AddText(testvars.TestDataID, testvars.TestText)
	assert.NoError(t, err, "error AddText calling")

	expectedResp := "Text " + testvars.TestDataID + " added successfully"

	assert.Equal(t, resp, expectedResp, "Output text didn't match expected")
}

func TestGettxtCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err, "error read request body")

		var jsonBody models.TextID
		err = json.Unmarshal(body, &jsonBody)
		assert.NoError(t, err, "error body unmarshalling")

		assert.Equal(t, jsonBody.DataID, testvars.TestDataID, "Request data id didn't match expected")

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, testvars.TestText)
	}))
	defer server.Close()

	baseURL = server.URL

	svc := newTextService()

	resp, err := svc.GetText(testvars.TestDataID)
	assert.NoError(t, err, "error GetText calling")

	assert.Equal(t, resp, testvars.TestText, "Output text didn't match expected")
}

func TestDeltxtCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err, "error read request body")

		var jsonBody models.TextID
		err = json.Unmarshal(body, &jsonBody)
		assert.NoError(t, err, "error body unmarshalling")

		assert.Equal(t, jsonBody.DataID, testvars.TestDataID, "Request data id didn't match expected")

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Text "+jsonBody.DataID+" deleted successfully")
	}))
	defer server.Close()

	baseURL = server.URL

	svc := newTextService()

	resp, err := svc.DelText(testvars.TestDataID)
	assert.NoError(t, err, "error DelText calling")

	expectedResp := "Text " + testvars.TestDataID + " deleted successfully"

	assert.Equal(t, resp, expectedResp, "Output text didn't match expected")
}
