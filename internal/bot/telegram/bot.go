package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"task-planner-bot/internal/database"
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
	Rep database.Repository
}

// NewBotHandler создает новый экземпляр BotHandler
func NewBotHandler(bot *tgbotapi.BotAPI, rep database.Repository) *BotHandler {
	return &BotHandler{
		Bot: bot,
		Rep: rep,
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
