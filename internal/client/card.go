package client

import (
	"fmt"

	"github.com/nu-kotov/GophKeeper/internal/client_service"
	"github.com/spf13/cobra"
)

func newCardService() *client_service.CardService {
	return &client_service.CardService{
		BaseURL:    baseURL,
		HTTPClient: httpClient,
		Key:        []byte(key),
		LoadCookie: LoadCookie,
		Encrypt:    Encrypt,
		Decrypt:    Decrypt,
	}
}

var addCardCmd = &cobra.Command{
	Use:   "addcard",
	Short: "Добавить данные банковской карты",
	Run: func(cmd *cobra.Command, args []string) {
		if id == "" || number == "" || expiry == "" || cvv == "" || holder == "" {
			fmt.Println("обязательные поля: --id, --number, --expiry, --cvv, --holder")
			return
		}

		svc := newCardService()

		resp, err := svc.AddCard(id, number, expiry, cvv, holder)
		if err != nil {
			fmt.Println("Ошибка:", err)
			return
		}

		fmt.Println(resp)
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

		svc := newCardService()

		card, err := svc.GetCard(id)
		if err != nil {
			fmt.Println("Ошибка:", err)
			return
		}

		fmt.Println("Полученные данные карты:")
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

		svc := newCardService()

		resp, err := svc.DelCard(id)
		if err != nil {
			fmt.Println("Ошибка:", err)
			return
		}

		fmt.Println(resp)
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
