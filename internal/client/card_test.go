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

func TestAddcardCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/card/add", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err, "error read request body")
		defer r.Body.Close()

		var jsonBody models.CardID
		err = json.Unmarshal(body, &jsonBody)
		assert.NoError(t, err, "error body unmarshalling")

		assert.Equal(t, jsonBody.DataID, testDataID, "Request data id didn't match expected")

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, "Card "+jsonBody.DataID+" added successfully")
	}))
	defer server.Close()

	baseURL = server.URL
	id = testDataID
	number = testCardNumber
	expiry = testCardExp
	cvv = testCVV
	holder = testUser
	key = testClientKey

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	AddCard(id, number, expiry, cvv, holder)

	w.Close()
	os.Stdout = old
	out, _ := io.ReadAll(r)

	expected_out := "Card " + testDataID + " added successfully\n"

	assert.Equal(t, string(out), expected_out, "Output text didn't match expected")
}

func TestGetcardCommand(t *testing.T) {
	id = testDataID
	card := models.CardPayload{Number: testCardNumber, Expiry: testCardExp, CVV: testCVV, Name: testUser}
	cardJSON, err := json.Marshal(card)
	assert.NoError(t, err, "error card data marshalling")

	encryptedTestCardData, err := Encrypt([]byte(key), string(cardJSON))
	assert.NoError(t, err, "error card data encryption")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/card/get", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err, "error read request body")

		var jsonBody models.CardID
		err = json.Unmarshal(body, &jsonBody)
		assert.NoError(t, err, "error body unmarshalling")

		assert.Equal(t, jsonBody.DataID, testDataID, "Request data id didn't match expected")

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, encryptedTestCardData)
	}))

	defer server.Close()

	baseURL = server.URL

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	GetCard(id)

	w.Close()
	os.Stdout = old
	out, _ := io.ReadAll(r)

	expectedOut := "Полученные данные карты:\nНомер карты: 4433062851071851\nСрок: 2/2027\nДержатель: test_login\nCVV: 333\n"
	assert.Equal(t, string(out), expectedOut, "Output text didn't match expected")
}

func TestDelcardCommand(t *testing.T) {
	id = testDataID

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/card/delete", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err, "error read request body")

		var jsonBody models.CardID
		err = json.Unmarshal(body, &jsonBody)
		assert.NoError(t, err, "error body unmarshalling")

		assert.Equal(t, jsonBody.DataID, testDataID, "Request data id didn't match expected")

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Card "+id+" deleted successfully")
	}))

	defer server.Close()

	baseURL = server.URL

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	DelCard(id)

	w.Close()
	os.Stdout = old
	out, _ := io.ReadAll(r)

	expectedOut := "Card " + testDataID + " deleted successfully\n"
	assert.Equal(t, string(out), expectedOut, "Output text didn't match expected")
}
