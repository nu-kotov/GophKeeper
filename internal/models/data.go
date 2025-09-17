package models

type PrivateData struct {
	DataName      string `json:"data_name"`
	Login         string `json:"login"`
	Password      []byte `json:"password"`
	Text          string `json:"text"`
	BinaryContent []byte `json:"binary_content"`
	Card          []byte `json:"card"`
}
