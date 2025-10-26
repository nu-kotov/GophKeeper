package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/nu-kotov/GophKeeper/internal/testvars"
	"github.com/stretchr/testify/assert"
)

func TestBuildJWTString_And_GetUserID(t *testing.T) {
	secret := testvars.TestSecret
	userID := testvars.TestUserID
	login := testvars.TestUser

	token, err := BuildJWTString(userID, login, time.Minute*5, secret)
	assert.NoError(t, err)
	assert.NotEmpty(t, token, "empty token")

	gotUserID, err := GetUserID(token, secret)
	assert.NoError(t, err)
	assert.Equal(t, userID, gotUserID, "userID must match")
}

func TestGetUserID_InvalidToken(t *testing.T) {
	secret := testvars.TestSecret

	_, err := GetUserID("invalid.token.string", secret)
	assert.Error(t, err, "should be an error if the token is invalid")
}

func TestGetUserID_WrongSecret(t *testing.T) {
	secret := testvars.TestSecret
	wrongSecret := "wrongsecret"
	userID := testvars.TestUserID
	login := testvars.TestUser

	token, err := BuildJWTString(userID, login, time.Minute*5, secret)
	assert.NoError(t, err)

	_, err = GetUserID(token, wrongSecret)
	assert.Error(t, err, "should be an error if the secret key is incorrect")
}

func TestGetUserID_ExpiredToken(t *testing.T) {
	secret := testvars.TestSecret
	userID := testvars.TestUserID
	login := testvars.TestUser

	token, err := BuildJWTString(userID, login, -time.Minute, secret)
	assert.NoError(t, err)

	_, err = GetUserID(token, secret)
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "token is expired") || strings.Contains(err.Error(), "expired"), "error must be related to the token expiration")
}
