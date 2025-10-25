package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexedwards/argon2id"
	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/nu-kotov/GophKeeper/internal/auth"
	"github.com/nu-kotov/GophKeeper/internal/config"
	"github.com/nu-kotov/GophKeeper/internal/dberrors"
	"github.com/nu-kotov/GophKeeper/internal/mocks"
	"github.com/nu-kotov/GophKeeper/internal/models"
	"github.com/nu-kotov/GophKeeper/internal/testvars"
	"github.com/stretchr/testify/assert"
)

func TestUsersHandler_RegisterUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockUsersStorage(ctrl)
	mockStorage.EXPECT().InsertUserData(gomock.Any(), gomock.Any()).Return(nil)

	r := chi.NewRouter()

	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	config.SecretKey = testvars.TestSecret

	NewUsersHandler(r, config, mockStorage)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/user/register", map[string]string{
		"login":    testvars.TestUser,
		"password": testvars.TestPassword,
	},
		false,
	)

	token := resp.Cookies()[0]
	userID, err := auth.GetUserID(token.Value, testvars.TestSecret)
	assert.NoError(t, err, "error user id getting")

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "error body reading")
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode, "Response status code didn't match expected")
	assert.NotEmpty(t, userID, "User ID is empty")
	assert.Equal(t, "User "+testvars.TestUser+" registered successfully", string(body), "Response text didn't match expected")
}

func TestUsersHandler_RegisterUser_AlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockUsersStorage(ctrl)
	mockStorage.EXPECT().InsertUserData(gomock.Any(), gomock.Any()).Return(dberrors.ErrConflict)

	r := chi.NewRouter()

	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	config.SecretKey = testvars.TestSecret

	NewUsersHandler(r, config, mockStorage)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/user/register", map[string]string{
		"login":    testvars.TestUser,
		"password": testvars.TestPassword,
	},
		true,
	)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, 409, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "User "+testvars.TestUser+" already exists", string(body), "Response text didn't match expected")
}

func TestUsersHandler_LoginUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	passwordHash, err := argon2id.CreateHash(testvars.TestPassword, argon2id.DefaultParams)
	assert.NoError(t, err, "error password hashing")

	var userData = models.UserData{
		UserID:   testvars.TestUserID,
		Login:    testvars.TestUser,
		Password: passwordHash,
	}

	mockStorage := mocks.NewMockUsersStorage(ctrl)
	mockStorage.EXPECT().SelectUserData(gomock.Any(), gomock.Any()).Return(&userData, nil)

	r := chi.NewRouter()
	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	config.SecretKey = testvars.TestSecret

	NewUsersHandler(r, config, mockStorage)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/user/login", map[string]string{
		"login":    testvars.TestUser,
		"password": testvars.TestPassword,
	},
		false)

	token := resp.Cookies()[0]
	userID, err := auth.GetUserID(token.Value, testvars.TestSecret)
	assert.NoError(t, err, "error user id getting")

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, testvars.TestUserID, userID, "Response user id didn't match expected")
	assert.Equal(t, "User "+testvars.TestUser+" authorized", string(body), "Response text didn't match expected")
}

func TestUsersHandler_LoginUser_WrongPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	passwordHash, err := argon2id.CreateHash(testvars.TestPassword, argon2id.DefaultParams)
	assert.NoError(t, err, "error password hashing")

	var userData = models.UserData{
		UserID:   testvars.TestUserID,
		Login:    testvars.TestUser,
		Password: passwordHash,
	}

	mockStorage := mocks.NewMockUsersStorage(ctrl)
	mockStorage.EXPECT().SelectUserData(gomock.Any(), gomock.Any()).Return(&userData, nil)

	r := chi.NewRouter()
	config, err := config.NewConfig()
	assert.NoError(t, err, "error config init")

	config.SecretKey = testvars.TestSecret

	NewUsersHandler(r, config, mockStorage)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := makeRequest(t, ts, http.MethodPost, "/api/user/login", map[string]string{
		"login":    testvars.TestUser,
		"password": testvars.TestPassword + "wrong",
	},
		false)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Response status code didn't match expected")
	assert.Equal(t, "Uncorrect passwort", string(body), "Response text didn't match expected")
}
