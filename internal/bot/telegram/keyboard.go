package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"task-planner-bot/internal/consts"
)

func mainKeyboard(chatID int64) tgbotapi.Chattable {
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(consts.KeyNewTask, consts.NewTask)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(consts.KeyTasks, consts.Tasks)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(consts.KeySettings, consts.Settings)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(consts.KeyReport, consts.Report)),
	)

	// Отправка приветственного сообщения с клавиатурой
	msg := tgbotapi.NewMessage(chatID, "Выберите действие:")
	msg.ReplyMarkup = inlineKeyboard

	return msg
}

func settingsKeyboard(chatID int64) tgbotapi.Chattable {
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(consts.KeyNotify, consts.Notify)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(consts.KeyNotifyUntil, consts.NotifyUntil)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(consts.KeyCostOfTasks, consts.WorthOfTasks)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(consts.KeyBack, consts.Back)),
	)

	// Отправка приветственного сообщения с клавиатурой
	msg := tgbotapi.NewMessage(chatID, "Выберите настройку:")
	msg.ReplyMarkup = inlineKeyboard

	return msg
}

func keyboardState(chatID int64, text string) tgbotapi.Chattable {
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(consts.KeyEnable, consts.Enable)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(consts.KeyDisable, consts.Disable)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(consts.KeyCancel, consts.Back)),
	)

	// Отправка приветственного сообщения с клавиатурой
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = inlineKeyboard

	return msg
}

func keyboardBack(chatID int64, text string) tgbotapi.Chattable {
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(consts.KeyCancel, consts.Back)),
	)

	// Отправка приветственного сообщения с клавиатурой
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = inlineKeyboard

	return msg
}
