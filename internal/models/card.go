package models

type CardData struct {
	DataID string `json:"data_id"`
	Card   []byte `json:"card"`
}
