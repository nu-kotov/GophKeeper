package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/nu-kotov/GophKeeper/internal/keeper_errors"
	"github.com/nu-kotov/GophKeeper/internal/models"
)

// UsersStorage - структура хранилища пользователей.
type UsersStorage struct {
	stor *DBStorage
}

// NewUsersStorage — конструктор хранилища пользователей.
func NewUsersStorage(stor *DBStorage) *UsersStorage {
	return &UsersStorage{stor: stor}
}

// InsertUserData - вставка пользователя в бд.
func (usrs *UsersStorage) InsertUserData(ctx context.Context, data *models.UserData) error {

	sql := `INSERT INTO users (login, password) VALUES ($1, $2);`

	tx, err := usrs.stor.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		sql,
		data.Login,
		data.Password,
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

// SelectUserData - получение пользователя из бд.
func (usrs *UsersStorage) SelectUserData(ctx context.Context, data *models.UserData) (*models.UserData, error) {

	var userData models.UserData

	sql := `SELECT user_id, password FROM users WHERE login = $1;`

	row := usrs.stor.db.QueryRowContext(
		ctx,
		sql,
		data.Login,
	)

	err := row.Scan(&userData.UserID, &userData.Password)
	if err != nil {
		return nil, err
	}

	return &userData, nil
}

// Close - закрывает соединение с бд.
func (usrs *UsersStorage) Close() error {
	return usrs.stor.db.Close()
}
