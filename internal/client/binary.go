package client

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

var getBinaryCmd = &cobra.Command{
	Use:   "getbinary",
	Short: "Получить бинарные данные",
	Run: func(cmd *cobra.Command, args []string) {

		if id == "" {
			fmt.Println("обязательное поле: --id")
			return
		}

		cookie, err := LoadCookie()
		if err != nil {
			fmt.Println("Error:", err.Error())
			return
		}

		req, err := http.NewRequest("POST", baseURL+"/api/binary/get", strings.NewReader(id))
		if err != nil {
			fmt.Println("Error:", err.Error())
			return
		}

		req.AddCookie(cookie)
		req.Header.Add("Content-Type", "text/plain")

		resp, err := httpClient.Do(req)
		if err != nil {
			fmt.Println("Error:", err.Error())
			return
		}
		defer resp.Body.Close()

		if filePath == "" {
			filePath = filepath.Base(id)
		}

		outFile, err := os.Create(filePath)
		if err != nil {
			fmt.Println("Error:", err.Error())
			return
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, resp.Body)
		if err != nil {
			fmt.Println("Error:", err.Error())
			return
		}

		fmt.Println("Файл успешно загружен")
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

		cookie, err := LoadCookie()
		if err != nil {
			fmt.Println("Error:", err.Error())
			return
		}

		req, err := http.NewRequest("POST", baseURL+"/api/binary/delete", strings.NewReader(id))
		if err != nil {
			fmt.Println("Error:", err.Error())
			return
		}

		req.AddCookie(cookie)
		req.Header.Add("Content-Type", "text/plain")

		resp, err := httpClient.Do(req)
		if err != nil {
			fmt.Println("Error:", err.Error())
			return
		}
		defer resp.Body.Close()

		fmt.Println("Файл успешно удален")
	},
}

func init() {
	addBinaryCmd.Flags().StringVarP(&filePath, "path", "p", "", "Путь к файлу")
	getBinaryCmd.Flags().StringVarP(&id, "id", "i", "", "ID секрета")
	getBinaryCmd.Flags().StringVarP(&filePath, "path", "p", "", "Путь к файлу")
	delBinaryCmd.Flags().StringVarP(&id, "id", "i", "", "ID секрета")
}
