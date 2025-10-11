package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var login string
var password string

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login user command",
	Run: func(cmd *cobra.Command, args []string) {
		loginData := map[string]string{
			"login":    login,
			"password": password,
		}

		body, err := json.Marshal(loginData)
		if err != nil {
			fmt.Println("Error encoding login data:", err)
			return
		}

		resp, err := http.Post(baseURL+"/api/user/login", "application/json", bytes.NewReader(body))
		if err != nil {
			fmt.Println("Login request failed:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("Login failed with status:", resp.Status)
			return
		}

		for _, cookie := range resp.Cookies() {
			if cookie.Name == "token" {
				err := SaveCookie(cookie)
				if err != nil {
					fmt.Println("Failed to save session:", err)
				} else {
					fmt.Println("Login successful. Session saved.")
				}
				return
			}
		}

		fmt.Println("Login response did not contain a token cookie.")
	},
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register user command",
	Run: func(cmd *cobra.Command, args []string) {
		loginData := map[string]string{
			"login":    login,
			"password": password,
		}

		body, err := json.Marshal(loginData)
		if err != nil {
			fmt.Println("Error encoding login data:", err)
			return
		}

		resp, err := http.Post(baseURL+"/api/user/register", "application/json", bytes.NewReader(body))
		if err != nil {
			fmt.Println("Register request failed:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("Register failed with status:", resp.Status)
			return
		}

		for _, cookie := range resp.Cookies() {
			if cookie.Name == "token" {
				err := SaveCookie(cookie)
				if err != nil {
					fmt.Println("Failed to save session:", err)
				} else {
					fmt.Println("Register successful. Session saved.")
				}
				return
			}
		}

		fmt.Println("Register response did not contain a token cookie.")
	},
}

func init() {
	loginCmd.Flags().StringVarP(&login, "login", "l", "", "Login username")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "Login password")
	loginCmd.MarkFlagRequired("login")
	loginCmd.MarkFlagRequired("password")
	registerCmd.Flags().StringVarP(&login, "login", "l", "", "Register username")
	registerCmd.Flags().StringVarP(&password, "password", "p", "", "Register password")
	registerCmd.MarkFlagRequired("login")
	registerCmd.MarkFlagRequired("password")
}
