package custom_middleware

import (
	"context"
	"net/http"

	"github.com/backendengineerark/clients-api/pkg/customlogs"
	"github.com/google/uuid"
)

func LoggerWithCorrelationId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), customlogs.LoggerKey{}, customlogs.NewLogger(uuid.New().String()))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
