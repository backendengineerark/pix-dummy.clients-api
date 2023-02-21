package customlogs

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
)

type LoggerKey struct{}

func NewLogger(correlationId string) *log.Logger {
	prefix := fmt.Sprintf("CorrelationId=%s ", correlationId)
	logger := log.New(os.Stdout, prefix, log.Ldate|log.Ltime|log.Lmsgprefix)
	return logger
}

func AttachLoggerToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, LoggerKey{}, NewLogger(uuid.New().String()))
}

func ExtractLoggerFromContext(ctx context.Context) *log.Logger {
	return ctx.Value(LoggerKey{}).(*log.Logger)
}
