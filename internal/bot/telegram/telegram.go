package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"time"

	"task-planner-bot/internal/database"
)

type TaskData struct {
	date time.Time
	text string
}

type UserState struct {
	isInput  bool
	isNumber bool
	key      string
	state    string
	last     string
}

// BotHandler хранит экземпляр бота и подключение к базе данных
type BotHandler struct {
	Bot       *tgbotapi.BotAPI
	Rep       database.Repository
	userState map[int64]*UserState
	taskData  map[int64]*TaskData
}

// NewBotHandler создает новый экземпляр BotHandler
func NewBotHandler(bot *tgbotapi.BotAPI, rep database.Repository) *BotHandler {
	return &BotHandler{
		Bot:       bot,
		Rep:       rep,
		userState: make(map[int64]*UserState),
		taskData:  make(map[int64]*TaskData),
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

func (h *BotHandler) sandNewMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	h.sandMessage(msg)
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

func (h *BotHandler) changeState(chatID int64, isInput, isNumber bool, k, s string) {
	userState := *h.userState[chatID]

	userState.isInput = isInput
	userState.isNumber = isNumber
	userState.key = k
	userState.state = s

	h.userState[chatID] = &userState
}
