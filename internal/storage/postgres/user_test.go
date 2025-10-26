package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nu-kotov/GophKeeper/internal/dberrors"
	"github.com/nu-kotov/GophKeeper/internal/models"
	"github.com/nu-kotov/GophKeeper/internal/testvars"
	"github.com/stretchr/testify/assert"
)

func TestInsertUserData_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err, "error mock creation")
	defer db.Close()

	storage := &UsersStorage{Stor: &DBStorage{db: db}}
	ctx := context.Background()
	user := &models.UserData{Login: testvars.TestUser, Password: testvars.TestPassword}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO users`).
		WithArgs(testvars.TestUser, testvars.TestPassword).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = storage.InsertUserData(ctx, user)
	assert.NoError(t, err, "error data insertion")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInsertUserData_Conflict(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	storage := &UsersStorage{Stor: &DBStorage{db: db}}
	ctx := context.Background()
	user := &models.UserData{Login: testvars.TestUser, Password: testvars.TestPassword}

	pgErr := &pgconn.PgError{Code: "23505"}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO users`).
		WithArgs(testvars.TestUser, testvars.TestPassword).
		WillReturnError(pgErr)
	mock.ExpectRollback()

	err := storage.InsertUserData(ctx, user)
	assert.ErrorIs(t, err, dberrors.ErrConflict, "expected ErrConflict")
}

func TestInsertUserData_BeginError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	storage := &UsersStorage{Stor: &DBStorage{db: db}}
	ctx := context.Background()
	user := &models.UserData{Login: testvars.TestUser, Password: testvars.TestPassword}

	mock.ExpectBegin().WillReturnError(errors.New("begin failed"))

	err := storage.InsertUserData(ctx, user)
	if err == nil || err.Error() != "begin failed" {
		t.Errorf("expected begin failed error, got %v", err)
	}
}

func TestSelectUserData_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	storage := &UsersStorage{Stor: &DBStorage{db: db}}
	ctx := context.Background()

	input := &models.UserData{Login: testvars.TestUser}

	mock.ExpectQuery(`SELECT user_id, password FROM users WHERE login = \$1;`).
		WithArgs(testvars.TestUser).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "password"}).AddRow(testvars.TestUserID, testvars.TestEncriptedPass))

	user, err := storage.SelectUserData(ctx, input)
	assert.NoError(t, err, "data selection error")
	assert.Equal(t, user.UserID, testvars.TestUserID, "unexpected user_id")
	assert.Equal(t, user.Password, testvars.TestEncriptedPass, "unexpected user_id")
}

func TestSelectUserData_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	storage := &UsersStorage{Stor: &DBStorage{db: db}}
	ctx := context.Background()

	input := &models.UserData{Login: testvars.TestUser}

	mock.ExpectQuery(`SELECT user_id, password FROM users WHERE login = \$1;`).
		WithArgs(testvars.TestUser).
		WillReturnError(errors.New("db failed"))

	_, err := storage.SelectUserData(ctx, input)
	if err == nil || err.Error() != "db failed" {
		t.Errorf("expected db failed error, got %v", err)
	}
}
