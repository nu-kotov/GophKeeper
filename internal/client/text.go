package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/nu-kotov/GophKeeper/internal/models"
	"github.com/spf13/cobra"
)

var textID string
var text string

var sendTextCmd = &cobra.Command{
	Use:   "addtxt",
	Short: "Сохранить текстовую информацию",
	Run: func(cmd *cobra.Command, args []string) {
		msg, err := SendText(textID, text)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println(msg)
	},
}

func init() {
	sendTextCmd.Flags().StringVarP(&textID, "id", "i", "", "id секрета")
	sendTextCmd.Flags().StringVarP(&text, "txt", "t", "", "Текст")
	sendTextCmd.MarkFlagRequired("id")
	sendTextCmd.MarkFlagRequired("txt")
	rootCmd.AddCommand(sendTextCmd)
}

func SendText(textID string, txt string) (string, error) {
	data := models.TextData{
		DataID: textID,
		Text:   txt,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	cookie, err := LoadCookie()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", baseURL+"/api/text/add", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.AddCookie(cookie)
	req.Header.Add("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}
