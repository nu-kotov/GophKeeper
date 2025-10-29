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
	stor *DBStorage
}

// NewCardStorage — конструктор хранилища данных банковских карт.
func NewCardStorage(stor *DBStorage) *CardStorage {
	return &CardStorage{stor: stor}
}

// InsertCardData - вставка данных банковской карты в бд.
func (crd *CardStorage) InsertCardData(ctx context.Context, userID string, cardData *models.CardData) error {

	query := `INSERT INTO card (user_id, data_id, card_data) VALUES ($1, $2, $3);`

	tx, err := crd.stor.db.Begin()
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
func (crd *CardStorage) SelectCardData(ctx context.Context, userID string, dataID string) (string, error) {

	var card_data string

	query := `SELECT card_data FROM card WHERE user_id = $1 AND data_id = $2;`

	row := crd.stor.db.QueryRowContext(
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
func (crd *CardStorage) DeleteCardData(ctx context.Context, userID string, dataID string) error {

	query := `DELETE FROM card WHERE user_id = $1 AND data_id = $2;`

	_, err := crd.stor.db.ExecContext(
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
func (crd *CardStorage) Close() error {
	return crd.stor.db.Close()
}
