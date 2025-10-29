package client

import (
	"fmt"

	"github.com/nu-kotov/GophKeeper/internal/client_service"
	"github.com/spf13/cobra"
)

func newTextService() *client_service.TextService {
	return &client_service.TextService{
		BaseURL:    baseURL,
		HTTPClient: httpClient,
		Key:        []byte(key),
		LoadCookie: LoadCookie,
		Encrypt:    Encrypt,
		Decrypt:    Decrypt,
	}
}

var addTextCmd = &cobra.Command{
	Use:   "addtxt",
	Short: "Сохранить текстовую информацию",
	Run: func(cmd *cobra.Command, args []string) {

		svc := newTextService()

		resp, err := svc.AddText(id, text)
		if err != nil {
			fmt.Println("Ошибка:", err)
			return
		}

		fmt.Println(resp)
	},
}

var getTextCmd = &cobra.Command{
	Use:   "gettxt",
	Short: "Получить текстовую информацию",
	Run: func(cmd *cobra.Command, args []string) {

		svc := newTextService()

		txt, err := svc.GetText(id)
		if err != nil {
			fmt.Println("Ошибка:", err)
			return
		}

		fmt.Println(txt)
	},
}

var delTextCmd = &cobra.Command{
	Use:   "deltxt",
	Short: "Удалить текстовую информацию",
	Run: func(cmd *cobra.Command, args []string) {

		svc := newTextService()

		txt, err := svc.DelText(id)
		if err != nil {
			fmt.Println("Ошибка:", err)
			return
		}

		fmt.Println(txt)
	},
}

func init() {
	addTextCmd.Flags().StringVarP(&id, "id", "i", "", "id секрета")
	addTextCmd.Flags().StringVarP(&text, "txt", "t", "", "Текст")
	addTextCmd.MarkFlagRequired("id")
	addTextCmd.MarkFlagRequired("txt")
	getTextCmd.Flags().StringVarP(&id, "id", "i", "", "id секрета")
	getTextCmd.MarkFlagRequired("id")
	delTextCmd.Flags().StringVarP(&id, "id", "i", "", "id секрета")
	delTextCmd.MarkFlagRequired("id")
}
