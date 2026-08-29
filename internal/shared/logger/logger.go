package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(environment string) (*zap.Logger, error) {

	if environment == "development" {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return config.Build()
	}

	config := zap.NewProductionConfig()
	config.Encoding = "json"
	config.EncoderConfig.TimeKey = "timestamp"

	return config.Build()
}