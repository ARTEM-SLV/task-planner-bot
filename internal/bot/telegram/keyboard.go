package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"task-planner-bot/internal/consts"
)

func mainKeyboard(chatID int64) tgbotapi.Chattable {
	// Создаем inline-клавиатуру
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Создать задачу", "new_task")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Просмотреть задачи", "tasks")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Настройки", "settings")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Отчет", "report")),
	)

	// Отправка приветственного сообщения с клавиатурой
	msg := tgbotapi.NewMessage(chatID, "Выберите действие:")
	msg.ReplyMarkup = inlineKeyboard

	return msg
}

func settingsKeyboard(chatID int64) tgbotapi.Chattable {
	// Создаем inline-клавиатуру
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Уведомлять о задачах", consts.Notify)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Уведомлять за время", consts.NotifyUntil)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Учет ценности задач", consts.CostOfTasks)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Назад", "back")),
	)

	// Отправка приветственного сообщения с клавиатурой
	msg := tgbotapi.NewMessage(chatID, "Выберите настройку:")
	msg.ReplyMarkup = inlineKeyboard

	return msg
}

func keyboardEnableDisableK(chatID int64, value string) tgbotapi.Chattable {
	// Создаем inline-клавиатуру
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Включить", consts.Enable)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Отключить", consts.Disable)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Отмена", consts.Back)),
	)

	// Отправка приветственного сообщения с клавиатурой
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("В данный момент состояние '%s'", value))
	msg.ReplyMarkup = inlineKeyboard

	return msg
}
