package models

type Credentials struct {
	DataID   string `json:"data_id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type CredentialsID struct {
	DataID string `json:"data_id"`
}
