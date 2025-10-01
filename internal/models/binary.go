package models

type BinaryData struct {
	DataID        string `json:"data_id"`
	BinaryContent []byte `json:"binary_content"`
}

type BinaryID struct {
	DataID string `json:"data_id"`
}
