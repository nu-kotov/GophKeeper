package models

// CardData - структура запроса с данными банковской карты.
type CardData struct {
	DataID string `json:"data_id"`
	Card   string `json:"card"`
}

// CardID - структура запроса id банковской карты.
type CardID struct {
	DataID string `json:"data_id"`
}

// CardPayload - структура данных банковской карты.
type CardPayload struct {
	Number string `json:"number"`
	Expiry string `json:"expiry"`
	CVV    string `json:"cvv"`
	Name   string `json:"name"`
}
