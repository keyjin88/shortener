package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
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
		return err
	}
	// создаём новую конфигурацию логера
	cfg := zap.NewProductionConfig()
	// устанавливаем уровень
	cfg.Level = lvl
	//настраиваем формат даты
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	// создаём логер на основе конфигурации
	zl, err := cfg.Build()
	if err != nil {
		return err
	}

	// делаем регистратор SugaredLogger
	// устанавливаем синглтон
	Log = *zl.Sugar()
	return nil
}
