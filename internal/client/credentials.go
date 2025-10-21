package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

var addCredCmd = &cobra.Command{
	Use:   "addcred",
	Short: "Добавить логин и пароль",
	Run: func(cmd *cobra.Command, args []string) {
		if id == "" || login == "" || password == "" {
			fmt.Println("необходимо указать id, login и password")
			return
		}

		if len(key) != 32 {
			fmt.Println("ключ должен быть 32 байта (AES-256)")
			return
		}

		encryptedPassword, err := Encrypt([]byte(key), password)
		if err != nil {
			fmt.Println("ошибка при шифровании: ", err.Error())
			return
		}

		payload := map[string]string{
			"data_id":  id,
			"login":    login,
			"password": encryptedPassword,
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

		req, err := http.NewRequest("POST", baseURL+"/api/credentials/add", bytes.NewBuffer(body))
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

var getCredCmd = &cobra.Command{
	Use:   "getcred",
	Short: "Получить учетные данные по ID",
	Run: func(cmd *cobra.Command, args []string) {
		if id == "" {
			fmt.Println("нужно указать id через флаг --id")
			return
		}

		cookie, err := LoadCookie()
		if err != nil {
			fmt.Println("Error:", err.Error())
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

		req, err := http.NewRequest("POST", baseURL+"/api/credentials/get", bytes.NewBuffer(body))
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

		var result struct {
			ID       string `json:"data_id"`
			Login    string `json:"login"`
			Password string `json:"password"` // encrypted
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			fmt.Println("Error:", err)
			return
		}

		decryptedPassword, err := Decrypt([]byte(key), result.Password)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Полученные учетные данные:")
		fmt.Printf("ID: %s\n", result.ID)
		fmt.Printf("Login: %s\n", result.Login)
		fmt.Printf("Password: %s\n", decryptedPassword)
	},
}

var delCredCmd = &cobra.Command{
	Use:   "delcred",
	Short: "Удалить учетные данные по ID",
	Run: func(cmd *cobra.Command, args []string) {
		if id == "" {
			fmt.Println("нужно указать id через флаг --id")
			return
		}

		cookie, err := LoadCookie()
		if err != nil {
			fmt.Println("Error:", err.Error())
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

		req, err := http.NewRequest("POST", baseURL+"/api/credentials/delete", bytes.NewBuffer(body))
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
	addCredCmd.Flags().StringVarP(&id, "id", "i", "", "ID секрета")
	addCredCmd.Flags().StringVarP(&login, "login", "l", "", "Логин")
	addCredCmd.Flags().StringVarP(&password, "password", "p", "", "Пароль")
	getCredCmd.Flags().StringVarP(&id, "id", "i", "", "ID секрета")
	delCredCmd.Flags().StringVarP(&id, "id", "i", "", "ID секрета")
}
