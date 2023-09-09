package logger

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log будет доступен всему коду как синглтон.
// Никакой код навыка, кроме функции InitLogger, не должен модифицировать эту переменную.
// По умолчанию установлен no-op-логер, который не выводит никаких сообщений.
var Log zap.SugaredLogger

// Initialize инициализирует синглтон логера с необходимым уровнем логирования.
func Initialize(level string) error {
	// преобразуем текстовый уровень логирования в zap.AtomicLevel

	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return errors.Wrap(err, "failed to parse log level")
	}
	// создаём новую конфигурацию логера
	cfg := zap.NewProductionConfig()
	// устанавливаем уровень
	cfg.Level = lvl
	// настраиваем формат даты
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	// создаём логер на основе конфигурации
	zl, err := cfg.Build()
	if err != nil {
		return errors.Wrap(err, "failed to build logger")
	}
	// делаем регистратор SugaredLogger
	// устанавливаем синглтон
	Log = *zl.Sugar()
	return nil
}
