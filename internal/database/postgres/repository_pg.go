package postgres

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
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

	if !exists {
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

		log.Printf("Пользователь %s (ID: %d) успешно добавлен или обновлен.", username, userID)
	}

	return nil
}
