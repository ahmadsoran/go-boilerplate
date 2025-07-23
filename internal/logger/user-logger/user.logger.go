package userlogger

import (
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	userLoggers = make(map[uint]*zap.SugaredLogger)
	mu          sync.Mutex
)

func GetUserLogger(userID uint) *zap.SugaredLogger {
	mu.Lock()
	defer mu.Unlock()

	if logger, exists := userLoggers[userID]; exists {
		return logger
	}

	logDir := "logs"
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		zap.L().Sugar().Warnf("Failed to create logs directory: %v", err)
		return zap.L().Sugar()
	}

	// Convert uint to string for file name
	userIDStr := strconv.FormatUint(uint64(userID), 10)
	path := filepath.Join(logDir, "user_"+userIDStr+".log")

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		zap.L().Sugar().Warnf("Failed to create logger for user %d: %v", userID, err)
		return zap.L().Sugar()
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(file),
		zap.InfoLevel,
	)
	logger := zap.New(core, zap.AddCaller())
	sugar := logger.Sugar()
	userLoggers[userID] = sugar
	return sugar
}
