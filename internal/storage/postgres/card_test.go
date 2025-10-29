package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nu-kotov/GophKeeper/internal/keeper_errors"
	"github.com/nu-kotov/GophKeeper/internal/models"
	"github.com/nu-kotov/GophKeeper/internal/testvars"
	"github.com/stretchr/testify/assert"
)

func TestInsertCardData_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err, "error mock creation")
	defer db.Close()

	storage := &CardStorage{stor: &DBStorage{db: db}}
	ctx := context.Background()
	card := &models.CardData{
		DataID: testvars.TestDataID,
		Card:   testvars.TestCardNumber,
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO card`).
		WithArgs(testvars.TestUser, testvars.TestDataID, testvars.TestCardNumber).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = storage.InsertCardData(ctx, testvars.TestUser, card)
	assert.NoError(t, err, "insertion error")

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err, "expectations error")
}

func TestInsertCardData_ConflictError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	storage := &CardStorage{stor: &DBStorage{db: db}}
	ctx := context.Background()
	card := &models.CardData{DataID: testvars.TestDataID, Card: testvars.TestCardNumber}

	pgErr := &pgconn.PgError{Code: "23505"}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO card`).
		WithArgs(testvars.TestUser, testvars.TestDataID, testvars.TestCardNumber).
		WillReturnError(pgErr)
	mock.ExpectRollback()

	err := storage.InsertCardData(ctx, testvars.TestUser, card)
	assert.ErrorIs(t, err, keeper_errors.ErrConflict, "expected ErrConflict")
}

func TestSelectCardData_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	storage := &CardStorage{stor: &DBStorage{db: db}}
	ctx := context.Background()

	mock.ExpectQuery(`SELECT card_data FROM card WHERE user_id = \$1 AND data_id = \$2;`).
		WithArgs(testvars.TestUser, testvars.TestDataID).
		WillReturnRows(sqlmock.NewRows([]string{"card_data"}).AddRow(testvars.TestCardNumber))

	result, err := storage.SelectCardData(ctx, testvars.TestUser, testvars.TestDataID)
	assert.NoError(t, err, "selection data error")
	assert.Equal(t, result, testvars.TestCardNumber, "unexpected result data")
}

func TestSelectCardData_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	storage := &CardStorage{stor: &DBStorage{db: db}}
	ctx := context.Background()

	mock.ExpectQuery(`SELECT card_data FROM card WHERE user_id = \$1 AND data_id = \$2;`).
		WithArgs(testvars.TestUser, testvars.TestDataID).
		WillReturnError(sql.ErrNoRows)

	_, err := storage.SelectCardData(ctx, testvars.TestUser, testvars.TestDataID)
	assert.ErrorIs(t, err, keeper_errors.ErrNotFound, "expected ErrNotFound")
}

func TestDeleteCardData_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	storage := &CardStorage{stor: &DBStorage{db: db}}
	ctx := context.Background()

	mock.ExpectExec(`DELETE FROM card WHERE user_id = \$1 AND data_id = \$2;`).
		WithArgs(testvars.TestUser, testvars.TestDataID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := storage.DeleteCardData(ctx, testvars.TestUser, testvars.TestDataID)
	assert.NoError(t, err, "deletion data error")
}

func TestDeleteCardData_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	storage := &CardStorage{stor: &DBStorage{db: db}}
	ctx := context.Background()

	mock.ExpectExec(`DELETE FROM card WHERE user_id = \$1 AND data_id = \$2;`).
		WithArgs(testvars.TestUser, testvars.TestDataID).
		WillReturnError(errors.New("db failed"))

	err := storage.DeleteCardData(ctx, testvars.TestUser, testvars.TestDataID)
	assert.Error(t, err, "expected error, got nil")
}
