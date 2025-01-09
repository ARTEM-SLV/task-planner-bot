package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v4/pgxpool"
	"task-planner-bot/internal/logging"
)

var logger *logging.Logger

// SetLogger устанавливает логгер для бота
func SetLogger(l *logging.Logger) {
	logger = l
}

// BotHandler хранит экземпляр бота и подключение к базе данных
type BotHandler struct {
	Bot *tgbotapi.BotAPI
	DB  *pgxpool.Pool
}

// NewBotHandler создает новый экземпляр BotHandler
func NewBotHandler(bot *tgbotapi.BotAPI, db *pgxpool.Pool) *BotHandler {
	return &BotHandler{
		Bot: bot,
		DB:  db,
	}
}

// StartPolling запускает цикл обработки сообщений
func (h *BotHandler) StartPolling() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := h.Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // Игнорируем non-Message Updates
			continue
		}

		// Обработка команд
		switch update.Message.Command() {
		case "start":
			h.HandleStart(update)
		case "new_task":
			h.HandleNewTask(update)
		case "tasks":
			h.HandleTasks(update)
		case "settings":
			h.HandleSettings(update)
		case "report":
			h.HandleReport(update)
		case "help":
			h.HandleHelp(update)
		default:
			h.HandleUnknownCommand(update)
		}
	}
}
