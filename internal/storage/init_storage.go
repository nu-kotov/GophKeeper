package storage

import (
	"github.com/nu-kotov/GophKeeper/internal/config"
	"github.com/nu-kotov/GophKeeper/internal/storage/postgres"
)

func NewPgStorage(c *config.Config) (*postgres.DBStorage, error) {

	DBStorage, err := postgres.NewConnect(c.DatabaseConnection)
	if err != nil {
		return nil, err
	}

	return DBStorage, nil
}

func NewUsersStorage(pg *postgres.DBStorage) *postgres.UsersStorage {
	return &postgres.UsersStorage{Stor: pg}
}

func NewTextStorage(pg *postgres.DBStorage) *postgres.TextStorage {
	return &postgres.TextStorage{Stor: pg}
}
