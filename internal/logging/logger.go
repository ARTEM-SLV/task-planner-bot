package logging

import (
	"database/sql"
	"log"
	"time"
)

// Logger структура для записи логов
type Logger struct {
	DB *sql.DB
}

// NewLogger создает новый экземпляр логгера
func NewLogger(db *sql.DB) *Logger {
	return &Logger{DB: db}
}

// Log записывает лог в базу данных
func (l *Logger) Log(userID int, message string) {
	query := `INSERT INTO logs (user_id, date, text) VALUES ($1, $2, $3)`
	_, err := l.DB.Exec(query, userID, time.Now(), message)
	if err != nil {
		log.Printf("Ошибка записи лога: %v", err)
	}
	log.Printf("Лог сохранен: user_id=%d, message=%s", userID, message)
}
