package customlogs

import (
	"context"
	"fmt"
	"log"
	"os"
)

type LoggerKey struct{}

func NewLogger(correlationId string) *log.Logger {
	prefix := fmt.Sprintf("CorrelationId=%s ", correlationId)
	logger := log.New(os.Stdout, prefix, log.Ldate|log.Ltime|log.Lmsgprefix)
	return logger
}

func GetContextLogger(ctx context.Context) *log.Logger {
	return ctx.Value(LoggerKey{}).(*log.Logger)
}
