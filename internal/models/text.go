package models

type TextData struct {
	DataID string `json:"data_id"`
	Text   string `json:"text"`
}

type TextID struct {
	DataID string `json:"data_id"`
}
