package client

import (
	"fmt"

	"github.com/nu-kotov/GophKeeper/internal/client_service"
	"github.com/spf13/cobra"
)

func newCredentialsService() *client_service.CredentialsService {
	return &client_service.CredentialsService{
		BaseURL:    baseURL,
		HTTPClient: httpClient,
		Key:        []byte(key),
		LoadCookie: LoadCookie,
		Encrypt:    Encrypt,
		Decrypt:    Decrypt,
	}
}

var addCredCmd = &cobra.Command{
	Use:   "addcred",
	Short: "Добавить логин и пароль",
	Run: func(cmd *cobra.Command, args []string) {

		if id == "" || login == "" || password == "" {
			fmt.Println("необходимо указать id, login и password")
			return
		}

		svc := newCredentialsService()

		resp, err := svc.AddCreds(id, login, password)
		if err != nil {
			fmt.Println("Ошибка:", err)
			return
		}

		fmt.Println(resp)
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

		svc := newCredentialsService()

		cred, err := svc.GetCreds(id)
		if err != nil {
			fmt.Println("Ошибка:", err)
			return
		}

		fmt.Println("Полученные учетные данные:")
		fmt.Printf("ID: %s\n", cred.ID)
		fmt.Printf("Login: %s\n", cred.Login)
		fmt.Printf("Password: %s\n", cred.Password)
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

		svc := newCredentialsService()

		resp, err := svc.DelCreds(id)
		if err != nil {
			fmt.Println("Ошибка:", err)
			return
		}

		fmt.Println(resp)
	},
}

func init() {
	addCredCmd.Flags().StringVarP(&id, "id", "i", "", "ID секрета")
	addCredCmd.Flags().StringVarP(&login, "login", "l", "", "Логин")
	addCredCmd.Flags().StringVarP(&password, "password", "p", "", "Пароль")
	getCredCmd.Flags().StringVarP(&id, "id", "i", "", "ID секрета")
	delCredCmd.Flags().StringVarP(&id, "id", "i", "", "ID секрета")
}
