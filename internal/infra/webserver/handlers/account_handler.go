package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/backendengineerark/clients-api/internal/usecase"
	"github.com/backendengineerark/clients-api/pkg/conversions"
	"github.com/backendengineerark/clients-api/pkg/customlogs"
)

type AccountHandler struct {
	CreateAccountUseCase *usecase.CreateAccountUseCase
}

func NewAccountHandler(createAccountUseCase *usecase.CreateAccountUseCase) *AccountHandler {
	return &AccountHandler{
		CreateAccountUseCase: createAccountUseCase,
	}
}

func (ah *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	logger := customlogs.GetContextLogger(r.Context())

	var input usecase.AccountInputDTO

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		logger.Printf("Fail to convert body beacause %s", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	logger.Printf("Try to create account with input=%s", conversions.StructToJsonIgnoreErrors(r.Context(), input))

	output, errors, err := ah.CreateAccountUseCase.Execute(r.Context(), input)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(errors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors)
		return
	}

	logger.Printf("Success to process with result %s", conversions.StructToJsonIgnoreErrors(r.Context(), output))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}
