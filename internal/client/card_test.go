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

func TestAddcardCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/card/add", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err, "error read request body")
		defer r.Body.Close()

		var jsonBody models.CardID
		err = json.Unmarshal(body, &jsonBody)
		assert.NoError(t, err, "error body unmarshalling")

		assert.Equal(t, jsonBody.DataID, testvars.TestDataID, "Request data id didn't match expected")

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, "Card "+jsonBody.DataID+" added successfully")
	}))
	defer server.Close()

	baseURL = server.URL
	id = testvars.TestDataID
	number = testvars.TestCardNumber
	expiry = testvars.TestCardExp
	cvv = testvars.TestCVV
	holder = testvars.TestUser
	key = []byte(testvars.TestClientKey)

	svc := newCardService()

	resp, err := svc.AddCard(id, number, expiry, cvv, holder)
	assert.NoError(t, err, "error AddCard calling")

	expectedResp := "Card " + testvars.TestDataID + " added successfully"

	assert.Equal(t, resp, expectedResp, "Output text didn't match expected")
}

func TestGetcardCommand(t *testing.T) {
	id = testvars.TestDataID
	card := models.CardPayload{Number: testvars.TestCardNumber, Expiry: testvars.TestCardExp, CVV: testvars.TestCVV, Name: testvars.TestUser}
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

		assert.Equal(t, jsonBody.DataID, testvars.TestDataID, "Request data id didn't match expected")

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, encryptedTestCardData)
	}))

	defer server.Close()

	baseURL = server.URL

	svc := newCardService()

	resp, err := svc.GetCard(id)
	assert.NoError(t, err, "error GetCard calling")

	assert.Equal(t, resp.Number, testvars.TestCardNumber, "Output card number didn't match expected")
	assert.Equal(t, resp.Expiry, testvars.TestCardExp, "Output exp date didn't match expected")
	assert.Equal(t, resp.Name, testvars.TestUser, "Output holder name didn't match expected")
	assert.Equal(t, resp.CVV, testvars.TestCVV, "Output CVV didn't match expected")
}

func TestDelcardCommand(t *testing.T) {
	id = testvars.TestDataID

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/card/delete", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err, "error read request body")

		var jsonBody models.CardID
		err = json.Unmarshal(body, &jsonBody)
		assert.NoError(t, err, "error body unmarshalling")

		assert.Equal(t, jsonBody.DataID, testvars.TestDataID, "Request data id didn't match expected")

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Card "+id+" deleted successfully")
	}))

	defer server.Close()

	baseURL = server.URL

	svc := newCardService()

	resp, err := svc.DelCard(id)
	assert.NoError(t, err, "error GetCard calling")

	expectedResp := "Card " + testvars.TestDataID + " deleted successfully"
	assert.Equal(t, resp, expectedResp, "Output text didn't match expected")
}
