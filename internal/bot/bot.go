package bot

import (
	"task-planner-bot/internal/logging"
)

var logger *logging.Logger

// SetLogger устанавливает логгер для бота
func SetLogger(l *logging.Logger) {
	logger = l
}

type Bot interface {
	StartPolling()
}
