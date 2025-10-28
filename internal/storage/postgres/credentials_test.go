package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"

	"github.com/nu-kotov/GophKeeper/internal/keeper_errors"
	"github.com/nu-kotov/GophKeeper/internal/models"
	"github.com/nu-kotov/GophKeeper/internal/testvars"
)

func newMockStorage(t *testing.T) (*CredentialsStorage, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err, "error mock creation")

	stor := &CredentialsStorage{
		Stor: &DBStorage{db: db},
	}

	cleanup := func() { db.Close() }
	return stor, mock, cleanup
}

func TestInsertCredentialsData_Success(t *testing.T) {
	s, mock, cleanup := newMockStorage(t)
	defer cleanup()

	ctx := context.Background()
	userID := testvars.TestUserID
	cred := &models.Credentials{
		DataID:   testvars.TestDataID,
		Login:    testvars.TestUser,
		Password: testvars.TestPassword,
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO credentials`).
		WithArgs(userID, cred.DataID, cred.Login, cred.Password).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := s.InsertCredentialsData(ctx, userID, cred)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInsertCredentialsData_ConflictError(t *testing.T) {
	s, mock, cleanup := newMockStorage(t)
	defer cleanup()

	ctx := context.Background()
	userID := testvars.TestUserID
	cred := &models.Credentials{
		DataID:   testvars.TestDataID,
		Login:    testvars.TestUser,
		Password: testvars.TestPassword,
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO credentials`).
		WithArgs(userID, cred.DataID, cred.Login, cred.Password).
		WillReturnError(&pgconn.PgError{Code: "23505"})
	mock.ExpectRollback()

	err := s.InsertCredentialsData(ctx, userID, cred)
	assert.ErrorIs(t, err, keeper_errors.ErrConflict)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSelectCredentialsData_Success(t *testing.T) {
	s, mock, cleanup := newMockStorage(t)
	defer cleanup()

	ctx := context.Background()
	userID := testvars.TestUserID
	dataID := testvars.TestDataID

	rows := sqlmock.NewRows([]string{"data_id", "login", "password"}).
		AddRow(testvars.TestDataID, testvars.TestUser, testvars.TestPassword)

	mock.ExpectQuery(`SELECT data_id, login, password FROM credentials`).
		WithArgs(userID, dataID).
		WillReturnRows(rows)

	cred, err := s.SelectCredentialsData(ctx, userID, dataID)
	assert.NoError(t, err)
	assert.Equal(t, testvars.TestUser, cred.Login)
	assert.Equal(t, testvars.TestPassword, cred.Password)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSelectCredentialsData_NotFound(t *testing.T) {
	s, mock, cleanup := newMockStorage(t)
	defer cleanup()

	ctx := context.Background()
	userID := testvars.TestUserID
	dataID := testvars.TestDataID

	mock.ExpectQuery(`SELECT data_id, login, password FROM credentials`).
		WithArgs(userID, dataID).
		WillReturnError(sql.ErrNoRows)

	cred, err := s.SelectCredentialsData(ctx, userID, dataID)
	assert.Nil(t, cred)
	assert.ErrorIs(t, err, keeper_errors.ErrNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteCredentialsData_Success(t *testing.T) {
	s, mock, cleanup := newMockStorage(t)
	defer cleanup()

	ctx := context.Background()
	userID := testvars.TestUserID
	dataID := testvars.TestDataID

	mock.ExpectExec(`DELETE FROM credentials`).
		WithArgs(userID, dataID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := s.DeleteCredentialsData(ctx, userID, dataID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteCredentialsData_Error(t *testing.T) {
	s, mock, cleanup := newMockStorage(t)
	defer cleanup()

	ctx := context.Background()
	userID := testvars.TestUserID
	dataID := testvars.TestDataID

	mock.ExpectExec(`DELETE FROM credentials`).
		WithArgs(userID, dataID).
		WillReturnError(errors.New("db failure"))

	err := s.DeleteCredentialsData(ctx, userID, dataID)
	assert.EqualError(t, err, "db failure")
	assert.NoError(t, mock.ExpectationsWereMet())
}
