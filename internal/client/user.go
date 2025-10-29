package client

import (
	"fmt"

	"github.com/nu-kotov/GophKeeper/internal/client_service"
	"github.com/spf13/cobra"
)

func newUserService() *client_service.UserService {
	return &client_service.UserService{
		BaseURL:    baseURL,
		HTTPClient: httpClient,
		Key:        []byte(key),
		SaveCookie: SaveCookie,
		LoadCookie: LoadCookie,
		Encrypt:    Encrypt,
		Decrypt:    Decrypt,
	}
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login user command",
	Run: func(cmd *cobra.Command, args []string) {

		svc := newUserService()

		resp, err := svc.Login(login, password)
		if err != nil {
			fmt.Println("Ошибка:", err)
			return
		}

		fmt.Println(resp)
	},
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register user command",
	Run: func(cmd *cobra.Command, args []string) {
		svc := newUserService()

		resp, err := svc.Register(login, password)
		if err != nil {
			fmt.Println("Ошибка:", err)
			return
		}

		fmt.Println(resp)
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
