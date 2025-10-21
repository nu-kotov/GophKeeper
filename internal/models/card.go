package models

type CardData struct {
	DataID string `json:"data_id"`
	Card   string `json:"card"`
}

type CardID struct {
	DataID string `json:"data_id"`
}

type CardPayload struct {
	Number string `json:"number"`
	Expiry string `json:"expiry"`
	CVV    string `json:"cvv"`
	Name   string `json:"name"`
}
