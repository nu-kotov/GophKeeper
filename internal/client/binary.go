package client

import (
	"fmt"

	"github.com/nu-kotov/GophKeeper/internal/client_service"
	"github.com/spf13/cobra"
)

func newBinaryService() *client_service.BinaryService {
	return &client_service.BinaryService{
		BaseURL:    baseURL,
		HTTPClient: httpClient,
		Key:        []byte(key),
		LoadCookie: LoadCookie,
		Encrypt:    Encrypt,
		Decrypt:    Decrypt,
	}
}

var addBinaryCmd = &cobra.Command{
	Use:   "addbinary",
	Short: "Загрузить бинарные данные",
	Run: func(cmd *cobra.Command, args []string) {

		if filePath == "" {
			fmt.Println("нужно указать путь к файлу --path")
			return
		}

		svc := newBinaryService()

		resp, err := svc.AddBinary(id)
		if err != nil {
			fmt.Println("Ошибка:", err)
			return
		}

		fmt.Println(string(resp))
	},
}

var getBinaryCmd = &cobra.Command{
	Use:   "getbinary",
	Short: "Получить бинарные данные",
	Run: func(cmd *cobra.Command, args []string) {

		if id == "" {
			fmt.Println("обязательное поле: --id")
			return
		}

		svc := newBinaryService()

		err := svc.GetBinary(id, filePath)
		if err != nil {
			fmt.Println("Ошибка:", err)
			return
		}

		fmt.Println("Файл успешно скачан")
	},
}

var delBinaryCmd = &cobra.Command{
	Use:   "delbinary",
	Short: "Удалить бинарные данные",
	Run: func(cmd *cobra.Command, args []string) {

		if id == "" {
			fmt.Println("обязательное поле: --id")
			return
		}

		svc := newBinaryService()

		resp, err := svc.DelBinary(id)
		if err != nil {
			fmt.Println("Ошибка:", err)
			return
		}

		fmt.Println(string(resp))
	},
}

func init() {
	addBinaryCmd.Flags().StringVarP(&filePath, "path", "p", "", "Путь к файлу")
	getBinaryCmd.Flags().StringVarP(&id, "id", "i", "", "ID секрета")
	getBinaryCmd.Flags().StringVarP(&filePath, "path", "p", "", "Путь к файлу")
	delBinaryCmd.Flags().StringVarP(&id, "id", "i", "", "ID секрета")
}
