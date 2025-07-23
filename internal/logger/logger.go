// internal/logger/logger.go

package logger

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var SystemLog *zap.SugaredLogger
var APILog *zap.SugaredLogger

func Init() {
	// System logger
	sysFile, err := os.OpenFile("system.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Can't open system log file: %v", err)
	}
	sysWS := zapcore.AddSync(sysFile)
	sysCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		sysWS,
		zap.InfoLevel,
	)
	systemLogger := zap.New(sysCore, zap.AddCaller())
	SystemLog = systemLogger.Sugar()

	// API logger
	apiFile, err := os.OpenFile("api.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Can't open API log file: %v", err)
	}
	apiWS := zapcore.AddSync(apiFile)
	apiCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		apiWS,
		zap.InfoLevel,
	)
	apiLogger := zap.New(apiCore, zap.AddCaller())
	APILog = apiLogger.Sugar()
}
