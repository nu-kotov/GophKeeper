package keeper_errors

import "errors"

// ErrConflict - ошибка при вставке дубля в бд.
var ErrConflict = errors.New("data conflict")

// ErrNotFound - ошибка при отсутствии данных в бд.
var ErrNotFound = errors.New("data not found")

// ErrInvalisCredentials - ошибка при отсутствии/неверных кредах.
var ErrInvalidCredentials = errors.New("invalid credentials")
