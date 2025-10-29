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

func TestAddcredCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/credentials/add", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err, "error read request body")
		defer r.Body.Close()

		var jsonBody models.Credentials

		err = json.Unmarshal(body, &jsonBody)
		assert.NoError(t, err, "error body unmarshalling")

		assert.Equal(t, jsonBody.DataID, testvars.TestDataID, "Request data id didn't match expected")
		assert.Equal(t, jsonBody.Login, testvars.TestUser, "Request login didn't match expected")

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, "Credentials "+jsonBody.DataID+" added successfully")
	}))
	defer server.Close()

	baseURL = server.URL
	id = testvars.TestDataID
	login = testvars.TestUser
	password = testvars.TestPassword
	key = []byte(testvars.TestClientKey)

	svc := newCredentialsService()

	resp, err := svc.AddCreds(id, login, password)
	assert.NoError(t, err, "error AddCreds calling")

	expectedResp := "Credentials " + testvars.TestDataID + " added successfully"

	assert.Equal(t, resp, expectedResp, "output text didn't match expected")
}

func TestGetcredCommand(t *testing.T) {

	id = testvars.TestDataID
	login = testvars.TestUser
	password = testvars.TestPassword
	key = []byte(testvars.TestClientKey)

	encryptedPassword, err := Encrypt([]byte(key), password)
	assert.NoError(t, err, "error password encryption")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/credentials/get", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err, "error read request body")

		var jsonBody models.CredentialsID
		err = json.Unmarshal(body, &jsonBody)
		assert.NoError(t, err, "error body unmarshalling")

		assert.Equal(t, jsonBody.DataID, testvars.TestDataID, "Request data id didn't match expected")

		respTestCreds := models.Credentials{
			DataID:   testvars.TestDataID,
			Login:    testvars.TestUser,
			Password: encryptedPassword,
		}
		JSONResp, err := json.Marshal(respTestCreds)
		assert.NoError(t, err, "error credentials marshalling")

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write(JSONResp)
	}))

	defer server.Close()

	baseURL = server.URL

	svc := newCredentialsService()

	resp, err := svc.GetCreds(id)
	assert.NoError(t, err, "error GetCreds calling")

	assert.Equal(t, resp.ID, testvars.TestDataID, "Output data id didn't match expected")
	assert.Equal(t, resp.Login, testvars.TestUser, "Output login didn't match expected")
	assert.Equal(t, resp.Password, testvars.TestPassword, "Output password didn't match expected")
}

func TestDelcredCommand(t *testing.T) {
	id = testvars.TestDataID

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/credentials/delete", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err, "error read request body")

		var jsonBody models.CredentialsID
		err = json.Unmarshal(body, &jsonBody)
		assert.NoError(t, err, "error body unmarshalling")

		assert.Equal(t, jsonBody.DataID, testvars.TestDataID, "Request data id didn't match expected")

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Credentials "+id+" deleted successfully")
	}))

	defer server.Close()

	baseURL = server.URL

	svc := newCredentialsService()

	resp, err := svc.DelCreds(id)
	assert.NoError(t, err, "error GetCreds calling")

	expectedResp := "Credentials " + testvars.TestDataID + " deleted successfully"
	assert.Equal(t, resp, expectedResp, "Output text didn't match expected")
}
