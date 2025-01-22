package database

import "database/sql"

type Setting struct {
	ValueS sql.NullString // Строковое значение
	ValueI sql.NullInt32  // Целочисленное значение
	ValueB sql.NullBool   // Логическое значение
}

type Repository interface {
	AddUser(userID int64, username string, lastMsgID int) error
	GetSetting(userID int64, key string) (*Setting, error)
	SaveSetting(userID int64, key, value any) error
}
