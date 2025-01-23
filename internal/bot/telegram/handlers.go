package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"task-planner-bot/internal/consts"
)

func (h *BotHandler) HandleQuery(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID

	switch update.Message.Command() {
	case consts.Start:
		h.HandleStart(update)
	case consts.NewTask:
		h.HandleNewTask(chatID)
	case consts.Tasks:
		h.HandleTasks(chatID)
	case consts.Settings:
		h.HandleSettings(chatID)
	case consts.Report:
		h.HandleReport(chatID)
	case consts.Help:
		h.HandleHelp(chatID)
	default:
		userID := update.Message.From.ID
		text := update.Message.Text
		h.HandleUserRequests(chatID, userID, text)
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
	case consts.NewTask:
		h.HandleNewTask(chatID)
	case consts.Tasks:
		h.HandleTasks(chatID)
	case consts.Settings:
		h.HandleSettings(chatID)
	case consts.Report:
		h.HandleReport(chatID)
	case consts.Notify:
		h.HandleSettingState(chatID, userID, consts.Notify)
	case consts.NotifyUntil:
		h.HandleSettingNumber(chatID, userID, consts.NotifyUntil, consts.MsgEnterValueNotifyUntil)
	case consts.WorthOfTasks:
		h.HandleSettingState(chatID, userID, consts.WorthOfTasks)
	case consts.Enable:
		v := consts.Enable
		h.HandleEnableDisableNotify(chatID, userID, v)
	case consts.Disable:
		v := consts.Disable
		h.HandleEnableDisableNotify(chatID, userID, v)
	case consts.Back:
		h.HandleBack(chatID)
	default:
		h.HandleUserRequests(chatID, userID, data)
	}
}
