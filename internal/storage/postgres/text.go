package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/nu-kotov/GophKeeper/internal/models"
)

// ErrConflict - ошибка при вставке дубля в бд.
var ErrConflict = errors.New("data conflict")

type TextStorage struct {
	Stor *DBStorage
}

func (usrs *TextStorage) InsertTextData(ctx context.Context, userID string, textData *models.TextData) error {

	sql := `INSERT INTO text_data (user_id, data_id, text_data) VALUES ($1, $2, $3);`

	tx, err := usrs.Stor.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		sql,
		userID,
		textData.DataID,
		textData.Text,
	)

	if err != nil {
		tx.Rollback()

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return ErrConflict
		}

		return err
	}

	return tx.Commit()
}

func (usrs *TextStorage) SelectTextData(ctx context.Context, userID string, dataID string) (string, error) {

	var text string

	sql := `SELECT text_data FROM text_data WHERE user_id = $1 AND data_id = $2;`

	row := usrs.Stor.db.QueryRowContext(
		ctx,
		sql,
		userID,
		dataID,
	)

	err := row.Scan(&text)
	if err != nil {
		return "", err
	}

	return text, nil
}

func (usrs *TextStorage) DeleteTextData(ctx context.Context, userID string, dataID string) error {

	sql := `DELETE FROM text_data WHERE user_id = $1 AND data_id = $2;`

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
