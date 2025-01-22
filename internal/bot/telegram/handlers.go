package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"task-planner-bot/internal/consts"
)

func (h *BotHandler) HandleQuery(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID

	switch update.Message.Command() {
	case "start":
		h.HandleStart(update)
	case "new_task":
		h.HandleNewTask(chatID)
	case "tasks":
		h.HandleTasks(chatID)
	case "settings":
		h.HandleSettings(chatID)
	case "report":
		h.HandleReport(chatID)
	case "help":
		h.HandleHelp(chatID)
	default:
		h.HandleUnknownCommand(chatID)
	}
}

func (h *BotHandler) HandleCallbackQuery(update tgbotapi.Update) {
	callback := update.CallbackQuery
	data := callback.Data

	// Удаляем сообщение с кнопками чтобы пользователь не мог выбрать другое действие
	h.deleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)

	chatID := callback.From.ID
	userID := callback.From.ID

	switch data {
	case "new_task":
		h.HandleNewTask(chatID)
	case "tasks":
		h.HandleTasks(chatID)
	case "settings":
		h.HandleSettings(chatID)
	case "report":
		h.HandleReport(chatID)
	case consts.Notify:
		h.HandleSettingNotify(chatID, userID)
	case consts.Enable:
		v := consts.Enable
		h.HandleEnableDisableNotify(chatID, userID, v)
	case consts.Disable:
		v := consts.Disable
		h.HandleEnableDisableNotify(chatID, userID, v)
	case consts.Back:
		h.HandleBack(chatID)
	default:
		h.HandleUnknownCommand(chatID)
	}
}
