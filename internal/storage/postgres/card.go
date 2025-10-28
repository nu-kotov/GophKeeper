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

// CardStorage - структура хранилища данных банковских карт.
type CardStorage struct {
	Stor *DBStorage
}

// InsertCardData - вставка данных банковской карты в бд.
func (usrs *CardStorage) InsertCardData(ctx context.Context, userID string, cardData *models.CardData) error {

	query := `INSERT INTO card (user_id, data_id, card_data) VALUES ($1, $2, $3);`

	tx, err := usrs.Stor.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		query,
		userID,
		cardData.DataID,
		cardData.Card,
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

// SelectCardData - получение данных банковской карты в бд.
func (usrs *CardStorage) SelectCardData(ctx context.Context, userID string, dataID string) (string, error) {

	var card_data string

	query := `SELECT card_data FROM card WHERE user_id = $1 AND data_id = $2;`

	row := usrs.Stor.db.QueryRowContext(
		ctx,
		query,
		userID,
		dataID,
	)

	err := row.Scan(&card_data)
	if err != nil {

		if err == sql.ErrNoRows {
			return "", keeper_errors.ErrNotFound
		}

		return "", err
	}

	return card_data, nil
}

// DeleteCardData - удаление данных банковской карты в бд.
func (usrs *CardStorage) DeleteCardData(ctx context.Context, userID string, dataID string) error {

	query := `DELETE FROM card WHERE user_id = $1 AND data_id = $2;`

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
