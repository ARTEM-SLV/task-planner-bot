package postgres

import (
	"context"
	"errors"
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
        FROM tasks_bot.settings
        WHERE user_id = $1 AND key = $2;
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

// SaveSetting сохраняет настройку пользователя
func (r *RepositoryPg) SaveSetting(userID int64, key, value any) error {
	var (
		query string
		err   error
	)

	switch v := value.(type) {
	case string:
		// Если значение строка, сохраняем в поле value_s
		query = `
			INSERT INTO tasks_bot.settings (user_id, key, value_s, value_i, value_b)
			VALUES ($1, $2, $3, NULL, NULL)
			ON CONFLICT (user_id, key) DO UPDATE
			SET value_s = EXCLUDED.value_s, value_i = NULL, value_b = NULL
		`
		_, err = r.pool.Exec(context.Background(), query, userID, key, v)

	case int:
		// Если значение число, сохраняем в поле value_i
		query = `
			INSERT INTO tasks_bot.settings (user_id, key, value_s, value_i, value_b)
			VALUES ($1, $2, NULL, $3, NULL)
			ON CONFLICT (user_id, key) DO UPDATE
			SET value_s = NULL, value_i = EXCLUDED.value_i, value_b = NULL
		`
		_, err = r.pool.Exec(context.Background(), query, userID, key, v)

	case bool:
		// Если значение булево, сохраняем в поле value_b
		query = `
			INSERT INTO tasks_bot.settings (user_id, key, value_s, value_i, value_b)
			VALUES ($1, $2, NULL, NULL, $3)
			ON CONFLICT (user_id, key) DO UPDATE
			SET value_s = NULL, value_i = NULL, value_b = EXCLUDED.value_b
		`
		_, err = r.pool.Exec(context.Background(), query, userID, key, v)

	default:
		// Если тип значения не поддерживается
		return errors.New("unsupported value type")
	}

	// Проверяем ошибку выполнения запроса
	if err != nil {
		log.Printf("Ошибка сохранения настройки: %v", err)
		return err
	}

	log.Printf("Настройка сохранена: key=%s, value=%v", key, value)
	return nil
}
