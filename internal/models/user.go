package models

// UserData - структура данных пользователя.
type UserData struct {
	UserID   string `json:"user_id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
