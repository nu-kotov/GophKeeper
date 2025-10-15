package client

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var addBinaryCmd = &cobra.Command{
	Use:   "addbinary",
	Short: "Загрузить бинарные данные",
	Run: func(cmd *cobra.Command, args []string) {

		if filePath == "" {
			fmt.Println("нужно указать путь к файлу --path")
			return
		}

		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println("ошибка при открытии файла ", err)
			return
		}
		defer file.Close()

		cookie, err := LoadCookie()
		if err != nil {
			fmt.Println("Error:", err.Error())
			return
		}

		filename := filepath.Base(filePath)

		req, err := http.NewRequest("POST", baseURL+"/api/binary/add", file)
		if err != nil {
			fmt.Println("Error:", err.Error())
			return
		}

		req.AddCookie(cookie)
		req.Header.Add("X-Filename", filename)
		req.Header.Add("Content-Type", "application/octet-stream")

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

func init() {
	addBinaryCmd.Flags().StringVarP(&filePath, "path", "p", "", "Путь к файлу")
}
