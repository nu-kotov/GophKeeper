package client

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/nu-kotov/GophKeeper/internal/models"
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

		assert.Equal(t, jsonBody.DataID, testDataID, "Request data id didn't match expected")
		assert.Equal(t, jsonBody.Login, testUser, "Request login didn't match expected")

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, "Credentials "+jsonBody.DataID+" added successfully")
	}))
	defer server.Close()

	baseURL = server.URL
	id = testDataID
	login = testUser
	password = testPassword
	key = testClientKey

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	AddCreds(id, login, password)

	w.Close()
	os.Stdout = old
	out, _ := io.ReadAll(r)

	expected_out := "Credentials " + testDataID + " added successfully\n"

	assert.Equal(t, string(out), expected_out, "Output text didn't match expected")
}

func TestGetcredCommand(t *testing.T) {

	id = testDataID
	login = testUser
	password = testPassword
	key = testClientKey

	encryptedPassword, err := Encrypt([]byte(key), password)
	assert.NoError(t, err, "error password encryption")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/credentials/get", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err, "error read request body")

		var jsonBody models.CredentialsID
		err = json.Unmarshal(body, &jsonBody)
		assert.NoError(t, err, "error body unmarshalling")

		assert.Equal(t, jsonBody.DataID, testDataID, "Request data id didn't match expected")

		respTestCreds := models.Credentials{
			DataID:   testDataID,
			Login:    testUser,
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

	GetCreds(id)

	w.Close()
	os.Stdout = old
	out, _ := io.ReadAll(r)

	expectedOut := "Полученные учетные данные:\nID: test_data\nLogin: test_login\nPassword: qwerty1\n"
	assert.Equal(t, string(out), expectedOut, "Output text didn't match expected")
}

func TestDelcredCommand(t *testing.T) {
	id = testDataID

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/credentials/delete", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err, "error read request body")

		var jsonBody models.CredentialsID
		err = json.Unmarshal(body, &jsonBody)
		assert.NoError(t, err, "error body unmarshalling")

		assert.Equal(t, jsonBody.DataID, testDataID, "Request data id didn't match expected")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Credentials "+id+" deleted successfully")
	}))

	defer server.Close()

	baseURL = server.URL

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	DelCreds(id)

	w.Close()
	os.Stdout = old
	out, _ := io.ReadAll(r)

	expectedOut := "Credentials " + testDataID + " deleted successfully\n"
	assert.Equal(t, string(out), expectedOut, "Output text didn't match expected")
}
