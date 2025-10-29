package client

import (
	"encoding/gob"
	"net/http"
	"os"
	"path/filepath"
)

const sessionFile = ".cli_session"

func getSessionPath() string {
	dir, _ := os.UserHomeDir()
	return filepath.Join(dir, sessionFile)
}

// Сохраняет сессию в файл
func SaveCookie(cookie *http.Cookie) error {
	file, err := os.Create(getSessionPath())
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(cookie)
}

// Загружает сессию из файла
func LoadCookie() (*http.Cookie, error) {
	file, err := os.Open(getSessionPath())
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cookie http.Cookie
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&cookie)
	if err != nil {
		return nil, err
	}

	return &cookie, nil
}
