package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/nu-kotov/GophKeeper/internal/dberrors"
	"github.com/nu-kotov/GophKeeper/internal/models"
)

type CredentialsStorage struct {
	Stor *DBStorage
}

func (usrs *CredentialsStorage) InsertCredentialsData(ctx context.Context, userID string, credentials *models.Credentials) error {

	sql := `INSERT INTO credentials (user_id, data_id, login, password) VALUES ($1, $2, $3, $4);`

	tx, err := usrs.Stor.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		sql,
		userID,
		credentials.DataID,
		credentials.Login,
		credentials.Password,
	)

	if err != nil {
		tx.Rollback()

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return dberrors.ErrConflict
		}

		return err
	}

	return tx.Commit()
}

func (usrs *CredentialsStorage) SelectCredentialsData(ctx context.Context, userID string, dataID string) (*models.Credentials, error) {

	var cred models.Credentials

	sql := `SELECT data_id, login, password FROM credentials WHERE user_id = $1 AND data_id = $2;`

	row := usrs.Stor.db.QueryRowContext(
		ctx,
		sql,
		userID,
		dataID,
	)

	err := row.Scan(&cred.DataID, &cred.Login, &cred.Password)
	if err != nil {
		return nil, err
	}

	return &cred, nil
}

func (usrs *CredentialsStorage) DeleteCredentialsData(ctx context.Context, userID string, dataID string) error {

	sql := `DELETE FROM credentials WHERE user_id = $1 AND data_id = $2;`

	_, err := usrs.Stor.db.ExecContext(
		ctx,
		sql,
		userID,
		dataID,
	)

	if err != nil {
		return err
	}

	return nil
}
