package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(levelString string) (*zap.Logger, error) {
	var level zapcore.Level

	if err := level.UnmarshalText([]byte(levelString)); err != nil {
		return nil, err
	}

	cfg := zap.Config{
		Level: zap.NewAtomicLevelAt(level),
	}

	logger, err := cfg.Build()

	if err != nil {
		return nil, err
	}

	return logger, nil
}
