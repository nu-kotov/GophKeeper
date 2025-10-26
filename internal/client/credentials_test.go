package client

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
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
	key = testvars.TestClientKey

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	addCreds(id, login, password)

	w.Close()
	os.Stdout = old
	out, _ := io.ReadAll(r)

	expected_out := "Credentials " + testvars.TestDataID + " added successfully\n"

	assert.Equal(t, string(out), expected_out, "Output text didn't match expected")
}

func TestGetcredCommand(t *testing.T) {

	id = testvars.TestDataID
	login = testvars.TestUser
	password = testvars.TestPassword
	key = testvars.TestClientKey

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

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	getCreds(id)

	w.Close()
	os.Stdout = old
	out, _ := io.ReadAll(r)

	expectedOut := "Полученные учетные данные:\nID: test_data\nLogin: test_login\nPassword: qwerty1\n"
	assert.Equal(t, string(out), expectedOut, "Output text didn't match expected")
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

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	delCreds(id)

	w.Close()
	os.Stdout = old
	out, _ := io.ReadAll(r)

	expectedOut := "Credentials " + testvars.TestDataID + " deleted successfully\n"
	assert.Equal(t, string(out), expectedOut, "Output text didn't match expected")
}
