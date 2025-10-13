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

var addTextCmd = &cobra.Command{
	Use:   "addtxt",
	Short: "Сохранить текстовую информацию",
	Run: func(cmd *cobra.Command, args []string) {
		msg, err := AddText(id, text)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println(msg)
	},
}

var getTextCmd = &cobra.Command{
	Use:   "gettxt",
	Short: "Получить текстовую информацию",
	Run: func(cmd *cobra.Command, args []string) {
		txt, err := GetText(id)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println(txt)
	},
}

var delTextCmd = &cobra.Command{
	Use:   "deltxt",
	Short: "Удалить текстовую информацию",
	Run: func(cmd *cobra.Command, args []string) {
		txt, err := DelText(id)
		if err != nil {
			fmt.Println("Error:", err)
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

func AddText(textID string, txt string) (string, error) {
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

func GetText(textID string) (string, error) {
	data := models.TextID{
		DataID: textID,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	cookie, err := LoadCookie()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", baseURL+"/api/text/get", bytes.NewBuffer(body))
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

func DelText(textID string) (string, error) {
	data := models.TextID{
		DataID: textID,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	cookie, err := LoadCookie()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", baseURL+"/api/text/delete", bytes.NewBuffer(body))
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
