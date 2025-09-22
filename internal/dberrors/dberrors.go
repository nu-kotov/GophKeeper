package dberrors

import "errors"

// ErrConflict - ошибка при вставке дубля в бд.
var ErrConflict = errors.New("data conflict")
