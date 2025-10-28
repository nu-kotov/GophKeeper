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

// TextStorage - структура хранилища под текст.
type TextStorage struct {
	Stor *DBStorage
}

// InsertTextData - вставка текста в бд.
func (usrs *TextStorage) InsertTextData(ctx context.Context, userID string, textData *models.TextData) error {

	query := `INSERT INTO text_data (user_id, data_id, text_data) VALUES ($1, $2, $3);`

	tx, err := usrs.Stor.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		query,
		userID,
		textData.DataID,
		textData.Text,
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

// SelectTextData - получение текста из бд.
func (usrs *TextStorage) SelectTextData(ctx context.Context, userID string, dataID string) (string, error) {

	var text string

	query := `SELECT text_data FROM text_data WHERE user_id = $1 AND data_id = $2;`

	row := usrs.Stor.db.QueryRowContext(
		ctx,
		query,
		userID,
		dataID,
	)

	err := row.Scan(&text)
	if err != nil {

		if err == sql.ErrNoRows {
			return "", keeper_errors.ErrNotFound
		}

		return "", err
	}

	return text, nil
}

// DeleteTextData - удаление текста из бд.
func (usrs *TextStorage) DeleteTextData(ctx context.Context, userID string, dataID string) error {

	query := `DELETE FROM text_data WHERE user_id = $1 AND data_id = $2;`

	_, err := usrs.Stor.db.ExecContext(
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
