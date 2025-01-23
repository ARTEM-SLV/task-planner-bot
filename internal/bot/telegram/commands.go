package telegram

import (
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"time"

	"task-planner-bot/internal/consts"
	"task-planner-bot/internal/database"
)

// HandleStart обрабатывает команду /start
func (h *BotHandler) HandleStart(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	h.userState[chatID] = &UserState{}

	msg := tgbotapi.NewMessage(chatID, consts.MsgWelcome)
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
	msg := tgbotapi.NewMessage(chatID, consts.MsgInputTaskDate)
	h.sandMessage(msg)

	// Здесь добавим логику для создания задачи
	//userState := h.userState[chatID]
	//userState.isInput = true
	//h.userState[chatID] = userState
	h.changeState(chatID, true, false, "", consts.EnumTaskDate)

}

// HandleTasks обрабатывает команду /tasks
func (h *BotHandler) HandleTasks(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, consts.MsgPeriodViewTasks)
	h.sandMessage(msg)
}

// HandleSettings обрабатывает команду /settings
func (h *BotHandler) HandleSettings(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, consts.MsgSettings)
	h.sandMessage(msg)

	// Здесь добавим логику для управления настройками пользователя
	newMsg := settingsKeyboard(chatID)
	h.sandMessage(newMsg)
}

// HandleReport обрабатывает команду /report
func (h *BotHandler) HandleReport(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, consts.MsgPeriodViewReport)
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

func (h *BotHandler) HandleSettingState(chatID, userID int64, key string) {
	setting, err := h.Rep.GetSetting(userID, key)
	if err != nil {
		log.Printf("Ошибка получения данных: %v", err)
	}

	if setting == nil {
		setting = &database.Setting{
			ValueB: sql.NullBool{
				Bool: false,
			},
		}
	}

	var value string
	if setting.ValueB.Bool {
		value = "Включено"
	} else {
		value = "Отключено"
	}

	//userState := h.userState[chatID]
	//userState.key = key
	//h.userState[chatID] = userState
	h.changeState(chatID, false, false, key, consts.EnumSettings)

	newMsg := keyboardState(chatID, value)
	h.sandMessage(newMsg)
}

func (h *BotHandler) HandleSettingNumber(chatID, userID int64, key string, text string) {
	setting, err := h.Rep.GetSetting(userID, key)
	if err != nil {
		log.Printf("Ошибка получения данных: %v", err)
	}

	var value string
	if setting == nil {
		value = fmt.Sprintf("'Отключено'")
	} else {
		value = fmt.Sprintf("%v мин.", setting.ValueI.Int32)
	}

	text = fmt.Sprintf("%s\n Сейчас значение %s Введите новое значение", text, value)

	//userState := h.userState[chatID]
	//userState.isInput = true
	//userState.key = key
	//userState.isNumber = true
	//h.userState[chatID] = userState
	h.changeState(chatID, true, true, key, consts.EnumSettings)

	newMsg := keyboardBack(chatID, text)
	h.sandMessage(newMsg)
}

func (h *BotHandler) HandleEnableDisableNotify(chatID, userID int64, v string) {
	var valueB bool
	key := h.userState[chatID].key

	if v == consts.Enable {
		valueB = true
	} else if v == consts.Disable {
		valueB = false
	} else {
		log.Printf("Unknown value: %s", v)
		return
	}

	err := h.Rep.SaveSetting(userID, key, valueB)
	if err != nil {
		log.Printf("Ошибка сохранения настроки %s: %v", key, err)
	}

	h.HandleBack(chatID)
}

func (h *BotHandler) HandleUserInput(chatID, userID int64, v any) {
	key := h.userState[chatID].key

	switch h.userState[chatID].state {
	case consts.EnumSettings:
		err := h.Rep.SaveSetting(userID, key, v)
		if err != nil {
			log.Printf("Ошибка сохранения настроки %s: %v", key, err)
		}

		h.changeState(chatID, false, false, "", "")
		h.HandleBack(chatID)
	case consts.EnumTaskDate:
		date, err := time.Parse("2006-01-02 15:04", v.(string))
		if err == nil {
			h.taskData[chatID] = &TaskData{
				date: date,
			}

			h.sandNewMessage(chatID, consts.MsgInputTaskText)
			h.changeState(chatID, true, false, "", consts.EnumTaskText)
		} else {
			msg := tgbotapi.NewMessage(chatID, "Неверный формат даты. Попробуйте еще раз в формате ГГГГ-ММ-ДД ЧЧ:ММ.")
			h.sandMessage(msg)
		}
	case consts.EnumTaskText:
		settings, err := h.Rep.GetSetting(userID, consts.WorthOfTasks)
		if err == nil {
			if settings.ValueB.Bool {
				taskDate := *h.taskData[chatID]
				taskDate.text = v.(string)
				h.taskData[chatID] = &taskDate

				h.sandNewMessage(chatID, consts.MsgInputTaskWorth)
				h.changeState(chatID, true, true, "", consts.EnumTuskWorth)
			}
		} else {
			log.Printf("Ошибка получения настроки %s %v:", consts.WorthOfTasks, err)
			h.changeState(chatID, false, false, "", "")
			h.HandleBack(chatID)
		}
	case consts.EnumTuskWorth:
		err := h.Rep.SaveTask(chatID, h.taskData[chatID].date, h.taskData[chatID].text, v.(int))
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Ошибка сохранения задачи. Попробуйте позже.")
			h.sandMessage(msg)
			log.Printf("Ошибка сохранения задачи: %v", err)
			h.HandleBack(chatID)
			return

			msg = tgbotapi.NewMessage(chatID, consts.MsgTaskAdded)
			h.sandMessage(msg)

			h.changeState(chatID, false, false, "", "")
			h.HandleBack(chatID)

			delete(h.taskData, chatID)
		}
	}

	//h.HandleBack(chatID)
}

func (h *BotHandler) HandleBack(chatID int64) {
	var newMsg tgbotapi.Chattable

	switch h.userState[chatID].last {
	case consts.Settings:
		newMsg = settingsKeyboard(chatID)
	case consts.Root:
		newMsg = mainKeyboard(chatID)
	default:
		newMsg = mainKeyboard(chatID)
	}

	//newMsg := mainKeyboard(chatID)
	h.sandMessage(newMsg)
}

// HandleUserRequests обрабатывает все неизвестные запросы от пользователя
func (h *BotHandler) HandleUserRequests(chatID, userID int64, txt string) {
	state, exists := h.userState[chatID]
	if !exists {
		msg := tgbotapi.NewMessage(chatID, "Извините, я не знаю такой команды.")
		h.sandMessage(msg)

		return
	}

	if state.isInput {
		if state.isNumber {
			number, err := strconv.Atoi(txt)
			if err == nil {
				h.HandleUserInput(chatID, userID, number)
			} else {
				msg := tgbotapi.NewMessage(chatID, consts.MsgErrEnteringNumber)
				_, _ = h.Bot.Send(msg)
			}
		} else {
			h.HandleUserInput(chatID, userID, txt)
		}

		//userState := h.userState[chatID]
		//userState.isInput = false
		//userState.isNumber = false
		//h.userState[chatID] = userState
		//h.changeState(chatID, "", "", false, false, )

		return
	}

	msg := tgbotapi.NewMessage(chatID, "Извините, я не знаю такой команды.")
	h.sandMessage(msg)
}
