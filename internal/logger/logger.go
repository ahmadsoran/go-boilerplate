// internal/logger/logger.go
package logger

import (
	"log"

	"go.uber.org/zap"
)

var Log *zap.SugaredLogger

func Init() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Can't init zap logger: %v", err)
	}
	Log = logger.Sugar()
}
