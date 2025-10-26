package models

// TextData - структура запроса с текстом.
type TextData struct {
	DataID string `json:"data_id"`
	Text   string `json:"text"`
}

// TextID - структура запроса с id текста.
type TextID struct {
	DataID string `json:"data_id"`
}
