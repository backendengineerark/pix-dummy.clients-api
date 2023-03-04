package custom_middleware

import (
	"context"
	"net/http"

	"github.com/backendengineerark/clients-api/pkg/customlogs"
	"github.com/google/uuid"
)

func TracingRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idempontencyKey, nullableBefore := GetIdempontencyKey(r)
		logger := customlogs.NewLogger(idempontencyKey)
		ctx := context.WithValue(r.Context(), customlogs.LoggerKey{}, logger)

		if nullableBefore {
			logger.Printf("IdempotencyKey not send, set a default value=%s", idempontencyKey)
		} else {
			logger.Printf("IdempotencyKey received with value=%s", idempontencyKey)
		}

		w.Header().Add("IdempotencyKey", idempontencyKey)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetIdempontencyKey(r *http.Request) (string, bool) {
	idempotencyKey := r.Header.Get("IdempotencyKey")
	if idempotencyKey == "" {
		return uuid.New().String(), true
	}
	return idempotencyKey, false
}
