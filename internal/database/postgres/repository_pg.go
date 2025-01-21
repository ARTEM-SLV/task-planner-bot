package postgres

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"task-planner-bot/internal/database"
	"time"
)

type RepositoryPg struct {
	pool *pgxpool.Pool
}

func NewRepositoryPg(pool *pgxpool.Pool) *RepositoryPg {
	return &RepositoryPg{pool: pool}
}

// AddUser добавляет нового пользователя или обновляет данные существующего
func (r *RepositoryPg) AddUser(userID int64, username string, lastMsgID int) error {
	var exists bool

	// Шаг 1: Проверяем, существует ли пользователь
	checkQuery := `SELECT EXISTS (SELECT 1 FROM tasks_bot.users WHERE id = $1)`
	err := r.pool.QueryRow(context.Background(), checkQuery, userID).Scan(&exists)
	if err != nil {
		log.Printf("Ошибка при проверке существования пользователя: %v", err)
		return err
	}

	if exists {
		log.Printf("Пользователь %s (ID: %d) уже зарегистрирован.", username, userID)
	} else {
		// SQL-запрос для добавления или обновления пользователя
		query := `
        INSERT INTO tasks_bot.users (id, username, last_msg_id, date_registration)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (id) DO UPDATE
        SET username = EXCLUDED.username,
            last_msg_id = EXCLUDED.last_msg_id;
    `

		// Выполнение запроса
		_, err = r.pool.Exec(context.Background(), query, userID, username, lastMsgID, time.Now())
		if err != nil {
			log.Printf("Ошибка при добавлении пользователя: %v", err)
			return err
		}

		log.Printf("Пользователь %s (ID: %d) успешно добавлен.", username, userID)
	}

	return nil
}

// GetSetting получает значение настройки по ключу для указанного пользователя
func (r *RepositoryPg) GetSetting(userID int64, key string) (*database.Setting, error) {
	query := `
        SELECT value_s, value_i, value_b
        FROM settings
        WHERE user_id = $1 AND setting_key = $2;
    `
	setting := &database.Setting{}
	err := r.pool.QueryRow(context.Background(), query, userID, key).Scan(&setting.ValueS, &setting.ValueI, &setting.ValueB)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil // Настройка не найдена
		}
		return nil, err
	}
	return setting, nil
}
