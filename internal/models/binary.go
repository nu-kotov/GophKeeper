package models

// BinaryData - структура запроса с бинарными данными.
type BinaryData struct {
	DataID        string `json:"data_id"`
	BinaryContent []byte `json:"binary_content"`
}

// BinaryData - структура запроса с id бинарных данных.
type BinaryID struct {
	DataID string `json:"data_id"`
}
