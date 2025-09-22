package models

type BinaryContent struct {
	DataID        string `json:"data_id"`
	BinaryContent []byte `json:"binary_content"`
}
