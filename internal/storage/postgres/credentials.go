package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/nu-kotov/GophKeeper/internal/keeper_errors"
	"github.com/nu-kotov/GophKeeper/internal/models"
)

// CredentialsStorage - структура хранилища под креды.
type CredentialsStorage struct {
	stor *DBStorage
}

// NewCredentialsStorage — конструктор хранилища кредов.
func NewCredentialsStorage(stor *DBStorage) *CredentialsStorage {
	return &CredentialsStorage{stor: stor}
}

// InsertCredentialsData - вставка кредов в бд.
func (cs *CredentialsStorage) InsertCredentialsData(ctx context.Context, userID string, credentials *models.Credentials) error {

	query := `INSERT INTO credentials (user_id, data_id, login, password) VALUES ($1, $2, $3, $4);`

	tx, err := cs.stor.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		query,
		userID,
		credentials.DataID,
		credentials.Login,
		credentials.Password,
	)

	if err != nil {
		tx.Rollback()

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return keeper_errors.ErrConflict
		}

		return err
	}

	return tx.Commit()
}

// SelectCredentialsData - получение кредов из бд.
func (cs *CredentialsStorage) SelectCredentialsData(ctx context.Context, userID string, dataID string) (*models.Credentials, error) {

	var cred models.Credentials

	query := `SELECT data_id, login, password FROM credentials WHERE user_id = $1 AND data_id = $2;`

	row := cs.stor.db.QueryRowContext(
		ctx,
		query,
		userID,
		dataID,
	)

	err := row.Scan(&cred.DataID, &cred.Login, &cred.Password)
	if err != nil {

		if err == sql.ErrNoRows {
			return nil, keeper_errors.ErrNotFound
		}

		return nil, err
	}

	return &cred, nil
}

// DeleteCredentialsData - удаление кредов из бд.
func (cs *CredentialsStorage) DeleteCredentialsData(ctx context.Context, userID string, dataID string) error {

	query := `DELETE FROM credentials WHERE user_id = $1 AND data_id = $2;`

	_, err := cs.stor.db.ExecContext(
		ctx,
		query,
		userID,
		dataID,
	)

	if err != nil {
		return err
	}

	return nil
}

// Close - закрывает соединение с бд.
func (cs *CredentialsStorage) Close() error {
	return cs.stor.db.Close()
}
