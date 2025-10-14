package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/phedde/luhn-algorithm"
	"github.com/spf13/cobra"
)

type CardPayload struct {
	Number string `json:"number"`
	Expiry string `json:"expiry"`
	CVV    string `json:"cvv"`
	Name   string `json:"name"`
}

var addCardCmd = &cobra.Command{
	Use:   "addcard",
	Short: "Добавить данные банковской карты",
	Run: func(cmd *cobra.Command, args []string) {
		if id == "" || number == "" || expiry == "" || cvv == "" || holder == "" {
			fmt.Println("обязательные поля: --id, --number, --expiry, --cvv, --holder")
			return
		}

		if len(key) != 32 {
			fmt.Println("ключ должен быть 32 байта (AES-256)")
			return
		}

		if err := validateCardData(number, expiry, cvv, holder); err != nil {
			fmt.Println("ошибка валидации данных : ", err.Error())
			return
		}

		card := CardPayload{Number: number, Expiry: expiry, CVV: cvv, Name: holder}
		cardJSON, err := json.Marshal(card)
		if err != nil {
			fmt.Println("Error:", err.Error())
			return
		}

		encryptedCardData, err := encrypt([]byte(key), string(cardJSON))
		if err != nil {
			fmt.Println("ошибка при шифровании: ", err.Error())
			return
		}

		payload := map[string]string{
			"data_id": id,
			"card":    encryptedCardData,
		}

		body, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("Error:", err.Error())
			return
		}

		cookie, err := LoadCookie()
		if err != nil {
			fmt.Println("Error:", err.Error())
			return
		}

		req, err := http.NewRequest("POST", baseURL+"/api/card/add", bytes.NewBuffer(body))
		if err != nil {
			fmt.Println("Error:", err.Error())
			return
		}
		req.AddCookie(cookie)
		req.Header.Add("Content-Type", "application/json")

		resp, err := httpClient.Do(req)
		if err != nil {
			fmt.Println("Error:", err.Error())
			return
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println(string(respBody))
	},
}

var getCardCmd = &cobra.Command{
	Use:   "getcard",
	Short: "Получить данные банковской карты",
	Run: func(cmd *cobra.Command, args []string) {
		if id == "" {
			fmt.Println("обязательное поле: --id")
			return
		}

		if len(key) != 32 {
			fmt.Println("ключ должен быть 32 байта (AES-256)")
			return
		}

		payload := map[string]string{
			"data_id": id,
		}

		body, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("Error:", err.Error())
			return
		}

		cookie, err := LoadCookie()
		if err != nil {
			fmt.Println("Error:", err.Error())
			return
		}

		req, err := http.NewRequest("POST", baseURL+"/api/card/get", bytes.NewBuffer(body))
		if err != nil {
			fmt.Println("Error:", err.Error())
			return
		}
		req.AddCookie(cookie)
		req.Header.Add("Content-Type", "application/json")

		resp, err := httpClient.Do(req)
		if err != nil {
			fmt.Println("Error:", err.Error())
			return
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		decryptedCardData, err := decrypt([]byte(key), string(respBody))
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		var card CardPayload
		err = json.Unmarshal([]byte(decryptedCardData), &card)
		if err != nil {
			fmt.Println("Error:", err.Error())
			return
		}

		fmt.Println("Полученные учетные данные:")
		fmt.Printf("Номер карты: %s\n", card.Number)
		fmt.Printf("Срок: %s\n", card.Expiry)
		fmt.Printf("Держатель: %s\n", card.Name)
		fmt.Printf("CVV: %s\n", card.CVV)
	},
}

var delCardCmd = &cobra.Command{
	Use:   "delcard",
	Short: "Удалить данные банковской карты",
	Run: func(cmd *cobra.Command, args []string) {
		if id == "" {
			fmt.Println("обязательное поле: --id")
			return
		}

		payload := map[string]string{
			"data_id": id,
		}

		body, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("Error:", err.Error())
			return
		}

		cookie, err := LoadCookie()
		if err != nil {
			fmt.Println("Error:", err.Error())
			return
		}

		req, err := http.NewRequest("POST", baseURL+"/api/card/delete", bytes.NewBuffer(body))
		if err != nil {
			fmt.Println("Error:", err.Error())
			return
		}
		req.AddCookie(cookie)
		req.Header.Add("Content-Type", "application/json")

		resp, err := httpClient.Do(req)
		if err != nil {
			fmt.Println("Error:", err.Error())
			return
		}
		defer resp.Body.Close()

		fmt.Println("Секрет удален")
	},
}

func init() {
	addCardCmd.Flags().StringVarP(&id, "id", "i", "", "ID секрета")
	addCardCmd.Flags().StringVarP(&number, "number", "n", "", "Номер карты")
	addCardCmd.Flags().StringVarP(&expiry, "expiry", "e", "", "Срок (MM/YY or MM/YYYY)")
	addCardCmd.Flags().StringVarP(&cvv, "cvv", "c", "", "CVV (3 digits)")
	addCardCmd.Flags().StringVarP(&holder, "holder", "u", "", "Имя держателя карты")
	getCardCmd.Flags().StringVarP(&id, "id", "i", "", "ID секрета")
	delCardCmd.Flags().StringVarP(&id, "id", "i", "", "ID секрета")
}

func validateCardData(number, expiry, cvv, holder string) error {
	// Проверка номера карты (Luhn + тип карты)
	intNumber, err := strconv.ParseInt(number, 10, 64)
	if err != nil {
		return err
	}
	if isValid := luhn.IsValid(intNumber); !isValid {
		return errors.New("invalid card number")
	}

	// CVV: ровно 3 цифры
	matched, _ := regexp.MatchString(`^\d{3}$`, cvv)
	if !matched {
		return errors.New("CVV must be exactly 3 digits")
	}

	// Expiry: не просрочен
	month, year, err := parseExpiryDate(expiry)
	if err != nil {
		return err
	}
	if isExpired(month, year) {
		return errors.New("card is expired")
	}

	if strings.TrimSpace(holder) == "" {
		return errors.New("holder is required")
	}

	return nil
}

func parseExpiryDate(expiry string) (int, int, error) {
	expiry = strings.TrimSpace(expiry)
	parts := strings.Split(expiry, "/")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("expiry must be in format MM/YY or MM/YYYY")
	}

	month, err := strconv.Atoi(parts[0])
	if err != nil || month < 1 || month > 12 {
		return 0, 0, fmt.Errorf("invalid month")
	}

	year, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid year")
	}

	if year < 100 {
		year += 2000
	}

	return month, year, nil
}

func isExpired(month, year int) bool {
	now := time.Now()
	exp := time.Date(year, time.Month(month)+1, 0, 23, 59, 59, 0, time.UTC)
	return now.After(exp)
}
