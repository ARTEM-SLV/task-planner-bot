package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
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
		if update.CallbackQuery != nil {
			// Обрабатываем callback-запрос (если нажали на кнопку)
			h.HandleCallbackQuery(update)
		} else if update.Message != nil {
			// Обрабатываем обычные сообщения
			h.HandleQuery(update)
		}
	}
}

func (h *BotHandler) sandMessage(msg tgbotapi.Chattable) {
	_, err := h.Bot.Send(msg)
	if err != nil {
		log.Printf("Ошибка отправки сообщения: %v", err)
	}
}

func (h *BotHandler) deleteMessage(chatID int64, msgID int) {
	deleteMsg := tgbotapi.NewDeleteMessage(chatID, msgID)
	_, err := h.Bot.Request(deleteMsg)
	if err != nil {
		log.Printf("Ошибка удаления сообщения: %v", err)
	} else {
		log.Printf("Сообщение удалено")
	}
}
