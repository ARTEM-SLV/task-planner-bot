package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// DB содержит пул подключений к базе данных
type DB struct {
	Pool *pgxpool.Pool
}

// InitDB инициализирует подключение к базе данных PostgreSQL
func InitDB() (*DB, error) {
	// Получаем параметры подключения из переменных окружения
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL не установлен")
	}

	// Создаем конфигурацию подключения с параметрами
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка разбора конфигурации подключения: %w", err)
	}

	// Устанавливаем таймауты и другие настройки
	config.MaxConns = 10                      // Максимальное количество подключений
	config.MinConns = 1                       // Минимальное количество подключений
	config.MaxConnLifetime = time.Hour        // Время жизни подключения
	config.MaxConnIdleTime = 30 * time.Minute // Время ожидания перед закрытием неактивного подключения

	// Подключаемся к базе данных
	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	log.Println("Подключение к базе данных установлено")
	return &DB{Pool: pool}, nil
}

// Close закрывает пул подключений к базе данных
func (db *DB) Close() {
	db.Pool.Close()
	log.Println("Подключение к базе данных закрыто")
}
