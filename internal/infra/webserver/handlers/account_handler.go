package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/backendengineerark/clients-api/internal/usecase"
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
	var input usecase.AccountInputDTO

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	output, errors, err := ah.CreateAccountUseCase.Execute(input)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(errors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}
