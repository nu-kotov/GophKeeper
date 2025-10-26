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

	msg, err := addText(testvars.TestDataID, testvars.TestText)
	assert.NoError(t, err, "error addTextCmd executing")

	expected_out := "Text " + testvars.TestDataID + " added successfully"

	assert.Equal(t, msg, expected_out, "Output text didn't match expected")
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

	msg, err := getText(testvars.TestDataID)
	assert.NoError(t, err, "error addTextCmd executing")

	expected_out := testvars.TestText

	assert.Equal(t, msg, expected_out, "Output text didn't match expected")
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

	msg, err := delText(testvars.TestDataID)
	assert.NoError(t, err, "error addTextCmd executing")

	expected_out := "Text " + testvars.TestDataID + " deleted successfully"

	assert.Equal(t, msg, expected_out, "Output text didn't match expected")
}
