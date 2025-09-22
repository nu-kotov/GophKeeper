package models

type Credentials struct {
	DataID   string `json:"data_id"`
	Login    string `json:"login"`
	Password []byte `json:"password"`
}
