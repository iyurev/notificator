package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New() (zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "time"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableCaller = true
	zaplogger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		return *zaplogger, err
	}
	return *zaplogger, nil
}

func ProjectName(projectName string) zap.Field {
	return zap.String("projectName", projectName)
}

func TelegramChatID(chatid int64) zap.Field {
	return zap.Int64("groupID", chatid)
}
