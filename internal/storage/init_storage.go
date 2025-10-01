package storage

import (
	"github.com/minio/minio-go/v7"
	"github.com/nu-kotov/GophKeeper/internal/config"
	"github.com/nu-kotov/GophKeeper/internal/storage/locminiio"
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

func NewCredentialsStorage(pg *postgres.DBStorage) *postgres.CredentialsStorage {
	return &postgres.CredentialsStorage{Stor: pg}
}

func NewCardStorage(pg *postgres.DBStorage) *postgres.CardStorage {
	return &postgres.CardStorage{Stor: pg}
}

func NewMiniIOStorage(c *config.Config) (*minio.Client, error) {

	MiniIOStorage, err := locminiio.NewMiniIOConnect(c.MiniIOConnection)
	if err != nil {
		return nil, err
	}

	return MiniIOStorage, nil
}

func NewBinaryStorage(miniio *minio.Client) *locminiio.BinaryStorage {
	return &locminiio.BinaryStorage{Stor: miniio}
}
