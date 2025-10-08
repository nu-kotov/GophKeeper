package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"

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

	u, _ := url.Parse(baseURL)

	jar, err := cookiejar.New(nil)
	if err != nil {
		return "", err
	}
	httpClient.Jar = jar
	httpClient.Jar.SetCookies(u, []*http.Cookie{
		{
			Name:  "token",
			Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjAyMjM0MTQsIlVzZXJJRCI6IjBlOTQ0OGQ3LTE2ZjQtNDEzMi04ZWMwLTMxZGE3NGM3ZWExOCIsIkxvZ2luIjoiIn0.UO6CsyHkjXYjkiUeByEwdurRtb4kRtCY9eLZVK5v2Wc",
			Path:  "/",
		},
	})

	resp, err := httpClient.Post(baseURL+"/api/text/add", "application/json", bytes.NewBuffer(body))
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
