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

type CardStorage struct {
	Stor *DBStorage
}

func (usrs *CardStorage) InsertCardData(ctx context.Context, userID string, cardData *models.CardData) error {

	sql := `INSERT INTO card (user_id, data_id, card_data) VALUES ($1, $2, $3);`

	tx, err := usrs.Stor.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		sql,
		userID,
		cardData.DataID,
		cardData.Card,
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

func (usrs *CardStorage) SelectCardData(ctx context.Context, userID string, dataID string) (string, error) {

	var card_data string

	sql := `SELECT card_data FROM card WHERE user_id = $1 AND data_id = $2;`

	row := usrs.Stor.db.QueryRowContext(
		ctx,
		sql,
		userID,
		dataID,
	)

	err := row.Scan(&card_data)
	if err != nil {
		return "", err
	}

	return card_data, nil
}

func (usrs *CardStorage) DeleteCardData(ctx context.Context, userID string, dataID string) error {

	sql := `DELETE FROM card WHERE user_id = $1 AND data_id = $2;`

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
