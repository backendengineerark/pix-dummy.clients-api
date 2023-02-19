package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/backendengineerark/clients-api/internal/usecase"
)

type AccountHandler struct {
}

func NewAccountHandler() *AccountHandler {
	return &AccountHandler{}
}

func (ah *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var input usecase.AccountInputDTO

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	fmt.Println(input)

	usecase := usecase.NewCreateAccountUseCase()

	output, errors := usecase.Execute(input)
	if len(errors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}
