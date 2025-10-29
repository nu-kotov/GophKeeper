package client

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	"github.com/caarlos0/env"
)

var (
	id       string
	text     string
	number   string
	expiry   string
	cvv      string
	holder   string
	login    string
	password string
	filePath string
	key      []byte
	baseURL  string
	cfg      Config
)

// Config конфигурция для клиента
type Config struct {
	EncryptionKey string `env:"ENCRYPTION_KEY" envDefault:"12345678901234567890123456789012"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8181"`
}

func init() {
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("Ошибка загрузки конфигурации: %v\n", err)
		os.Exit(1)
	}

	key = []byte(cfg.EncryptionKey)
	if len(key) != 32 {
		fmt.Println("Ошибка: длина ключа должна быть 32 байта (AES-256)")
		os.Exit(1)
	}

	fmt.Println("Конфигурация загружена")
	fmt.Printf("Base URL: %s\n", cfg.BaseURL)
}

// encrypt шифрует строку с использованием AES-GCM
func Encrypt(key []byte, plaintext string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decrypt выполняет расшифровку AES-GCM
func Decrypt(key []byte, encrypted string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(data) < aesGCM.NonceSize() {
		return "", errors.New("некорректные данные для расшифровки")
	}

	nonce := data[:aesGCM.NonceSize()]
	ciphertext := data[aesGCM.NonceSize():]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
