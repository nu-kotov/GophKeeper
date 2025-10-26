package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nu-kotov/GophKeeper/internal/dberrors"
	"github.com/nu-kotov/GophKeeper/internal/models"
	"github.com/nu-kotov/GophKeeper/internal/testvars"
	"github.com/stretchr/testify/assert"
)

func TestInsertTextData_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err, "error creating sqlmock")
	defer db.Close()

	storage := &TextStorage{Stor: &DBStorage{db: db}}
	ctx := context.Background()
	text := &models.TextData{DataID: testvars.TestDataID, Text: testvars.TestText}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO text_data`).
		WithArgs(testvars.TestUser, testvars.TestDataID, testvars.TestText).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = storage.InsertTextData(ctx, testvars.TestUser, text)
	assert.NoError(t, err, "insertion data error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInsertTextData_ConflictError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	storage := &TextStorage{Stor: &DBStorage{db: db}}
	ctx := context.Background()
	text := &models.TextData{DataID: testvars.TestDataID, Text: testvars.TestText}

	pgErr := &pgconn.PgError{Code: "23505"}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO text_data`).
		WithArgs(testvars.TestUser, testvars.TestDataID, testvars.TestText).
		WillReturnError(pgErr)
	mock.ExpectRollback()

	err := storage.InsertTextData(ctx, testvars.TestUser, text)
	assert.ErrorIs(t, err, dberrors.ErrConflict, "expected ErrConflict")
}

func TestSelectTextData_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	storage := &TextStorage{Stor: &DBStorage{db: db}}
	ctx := context.Background()

	mock.ExpectQuery(`SELECT text_data FROM text_data WHERE user_id = \$1 AND data_id = \$2;`).
		WithArgs(testvars.TestUser, testvars.TestDataID).
		WillReturnRows(sqlmock.NewRows([]string{"text_data"}).AddRow(testvars.TestText))

	result, err := storage.SelectTextData(ctx, testvars.TestUser, testvars.TestDataID)

	assert.NoError(t, err, "selection data error")
	assert.Equal(t, result, testvars.TestText, "unexpexted result")
}

func TestSelectTextData_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	storage := &TextStorage{Stor: &DBStorage{db: db}}
	ctx := context.Background()

	mock.ExpectQuery(`SELECT text_data FROM text_data WHERE user_id = \$1 AND data_id = \$2;`).
		WithArgs(testvars.TestUser, testvars.TestText).
		WillReturnError(sql.ErrNoRows)

	_, err := storage.SelectTextData(ctx, testvars.TestUser, testvars.TestText)
	assert.ErrorIs(t, err, dberrors.ErrNotFound, "expected ErrNotFound")
}

func TestDeleteTextData_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	storage := &TextStorage{Stor: &DBStorage{db: db}}
	ctx := context.Background()

	mock.ExpectExec(`DELETE FROM text_data WHERE user_id = \$1 AND data_id = \$2;`).
		WithArgs(testvars.TestUser, testvars.TestDataID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := storage.DeleteTextData(ctx, testvars.TestUser, testvars.TestDataID)
	assert.NoError(t, err, "deletion error")
}

func TestDeleteTextData_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	storage := &TextStorage{Stor: &DBStorage{db: db}}
	ctx := context.Background()

	mock.ExpectExec(`DELETE FROM text_data WHERE user_id = \$1 AND data_id = \$2;`).
		WithArgs(testvars.TestUser, testvars.TestDataID).
		WillReturnError(errors.New("db failed"))

	err := storage.DeleteTextData(ctx, testvars.TestUser, testvars.TestDataID)

	assert.Error(t, err, "expected error, got nil")
}
