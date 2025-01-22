package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5" // Telegram API
	"github.com/joho/godotenv"                                    // Для загрузки переменных из .env файла

	bot "task-planner-bot/internal/bot/telegram"
	"task-planner-bot/internal/database" // Модуль для работы с базой данных
	"task-planner-bot/internal/database/postgres"
)

func initLogger() {
	// Создаем файл для логирования с текущей датой
	logFile, err := os.OpenFile("bot.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Ошибка при создании файла логов: %v", err)
	}

	// Настраиваем логгер на запись в файл
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile) // Установка формата логов
	log.Println("Логирование началось")
}

func main() {
	// Загрузка переменных окружения из .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}

	initLogger()

	// Инициализация подключения к базе данных
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}
	defer db.Close()

	// Получаем токен для Telegram
	botToken := os.Getenv("TELEGRAM_TOKEN")
	if botToken == "" {
		log.Fatal("Токен Telegram не найден в переменных окружения")
	}

	// Инициализация бота
	botAPI, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Ошибка инициализации Telegram бота: %v", err)
	}

	botAPI.Debug = true
	log.Printf("Авторизован бот %s", botAPI.Self.UserName)

	rep := postgres.NewRepositoryPg(db.Pool)

	// Запуск обработчика команд бота
	botHandler := bot.NewBotHandler(botAPI, rep)
	botHandler.StartPolling()

	log.Println("Бот завершил работу")
}
