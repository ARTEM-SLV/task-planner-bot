package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleStart обрабатывает команду /start
func (h *BotHandler) HandleStart(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID

	msg := tgbotapi.NewMessage(chatID, "Привет! Я Task Planner Bot.\nЯ помогу вам организовать задачи.")
	h.sandMessage(msg)

	userID := update.Message.From.ID
	username := update.Message.From.UserName
	lastMsgID := update.Message.MessageID

	// Добавляем или обновляем пользователя
	err := h.Rep.AddUser(userID, username, lastMsgID)
	if err != nil {
		log.Printf("Ошибка добавления или обновления пользователя: %v", err)
		return
	}

	newMsg := mainKeyboard(chatID)
	h.sandMessage(newMsg)
}

// HandleNewTask обрабатывает команду /new_task
func (h *BotHandler) HandleNewTask(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Введите дату и текст задачи.")
	h.sandMessage(msg)
	// Здесь добавим логику для создания задачи
}

// HandleTasks обрабатывает команду /tasks
func (h *BotHandler) HandleTasks(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Выберите период для просмотра задач.")
	h.sandMessage(msg)
}

// HandleSettings обрабатывает команду /settings
func (h *BotHandler) HandleSettings(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Настройки: уведомления, период уведомлений, учет ценности задач.")
	h.sandMessage(msg)

	// Здесь добавим логику для управления настройками пользователя
	newMsg := settingsKeyboard(chatID)
	h.sandMessage(newMsg)
}

// HandleReport обрабатывает команду /report
func (h *BotHandler) HandleReport(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Выберите период для отчета по выполненным задачам.")
	h.sandMessage(msg)
	// Здесь добавим логику для вывода отчета
}

// HandleHelp обрабатывает команду /help
func (h *BotHandler) HandleHelp(chatID int64) {
	helpText := `
		/start - Начало работы с ботом
		/new_task - Создать новую задачу
		/tasks - Просмотр задач
		/settings - Настройки
		/report - Отчет по выполненным задачам
		/help - Справка по командам`
	msg := tgbotapi.NewMessage(chatID, helpText)
	h.sandMessage(msg)
}

// HandleUnknownCommand обрабатывает неизвестные команды
func (h *BotHandler) HandleUnknownCommand(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Извините, я не знаю такой команды.")
	h.sandMessage(msg)
}

func (h *BotHandler) HandleBack(chatID int64) {
	newMsg := mainKeyboard(chatID)
	h.sandMessage(newMsg)
}
