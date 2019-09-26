package logger

import (
	"github.com/IamStubborN/petstore/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(config config.Logger) (*zap.Logger, error) {
	logger, err := zap.Config{
		Encoding:    config.Encoding,
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths: config.OutputPaths,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}.Build()

	if err != nil {
		return nil, err
	}

	zap.ReplaceGlobals(logger)

	return logger, nil
}
