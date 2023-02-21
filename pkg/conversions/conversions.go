package conversions

import (
	"context"
	"encoding/json"

	"github.com/backendengineerark/clients-api/pkg/customlogs"
)

func StructToJsonIgnoreErrors(ctx context.Context, input interface{}) string {
	e, err := json.Marshal(input)
	if err != nil {
		logger := customlogs.ExtractLoggerFromContext(ctx)
		logger.Printf("Fail to convert struct to json because %s", err)
	}
	return string(e)
}
