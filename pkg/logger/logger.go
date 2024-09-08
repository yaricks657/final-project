package logger

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
)

type Logger struct {
	zerolog.Logger
}

// создание нового логгера
func New(logFilePath string) (*Logger, error) {
	// подготовка файла для записи
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	// настройка выводов логгера
	level := zerolog.InfoLevel
	multi := zerolog.MultiLevelWriter(file, os.Stdout)
	logger := zerolog.New(multi).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(level)

	return &Logger{logger}, nil
}

// LogError - логирование ошибок
func (l *Logger) LogError(msg string, err error) {
	l.Error().Err(err).Msg(msg)
}

// LogInfo - логирование информационных сообщений
func (l *Logger) LogInfo(msg ...string) {
	combinedMsg := strings.Join(msg, " ")
	l.Info().Msg(combinedMsg)
}

// LogWarn - логирование предупреждений
func (l *Logger) LogWarn(msg string) {
	l.Warn().Msg(msg)
}
