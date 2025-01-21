package database

type Setting struct {
	ValueS string // Строковое значение
	ValueI int64  // Целочисленное значение
	ValueB bool   // Логическое значение
}

type Repository interface {
	AddUser(userID int64, username string, lastMsgID int) error
	GetSetting(userID int64, key string) (*Setting, error)
}
