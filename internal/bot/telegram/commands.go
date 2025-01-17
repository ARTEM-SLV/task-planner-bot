package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleStart обрабатывает команду /start
func (h *BotHandler) HandleStart(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Я Task Planner Bot.\nЯ помогу вам организовать задачи.")
	h.Bot.Send(msg)

	userID := update.Message.From.ID
	username := update.Message.From.UserName
	lastMsgID := update.Message.MessageID

	// Добавляем или обновляем пользователя
	err := h.Rep.AddUser(userID, username, lastMsgID)
	if err != nil {
		log.Printf("Ошибка добавления или обновления пользователя: %v", err)
		return
	}
	// Здесь можно добавить логику для добавления пользователя в базу данных, если его нет
}

// HandleNewTask обрабатывает команду /new_task
func (h *BotHandler) HandleNewTask(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите дату и текст задачи.")
	h.Bot.Send(msg)
	// Здесь добавим логику для создания задачи
}

// HandleTasks обрабатывает команду /tasks
func (h *BotHandler) HandleTasks(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите период для просмотра задач.")
	h.Bot.Send(msg)
	// Здесь добавим логику для вывода задач по выбранному периоду
}

// HandleSettings обрабатывает команду /settings
func (h *BotHandler) HandleSettings(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Настройки: уведомления, период уведомлений, учет ценности задач.")
	h.Bot.Send(msg)
	// Здесь добавим логику для управления настройками пользователя
}

// HandleReport обрабатывает команду /report
func (h *BotHandler) HandleReport(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите период для отчета по выполненным задачам.")
	h.Bot.Send(msg)
	// Здесь добавим логику для вывода отчета
}

// HandleHelp обрабатывает команду /help
func (h *BotHandler) HandleHelp(update tgbotapi.Update) {
	helpText := `
		/start - Начало работы с ботом
		/new_task - Создать новую задачу
		/tasks - Просмотр задач
		/settings - Настройки
		/report - Отчет по выполненным задачам
		/help - Справка по командам`
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, helpText)
	h.Bot.Send(msg)
}

// HandleUnknownCommand обрабатывает неизвестные команды
func (h *BotHandler) HandleUnknownCommand(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Извините, я не знаю такой команды.")
	h.Bot.Send(msg)
}
