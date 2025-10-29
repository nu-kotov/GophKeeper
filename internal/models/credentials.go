package models

// Credentials - структура запроса с кредами.
type Credentials struct {
	DataID   string `json:"data_id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

// CredentialsID - структура запроса id кредов.
type CredentialsID struct {
	DataID string `json:"data_id"`
}
