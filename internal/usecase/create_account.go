package usecase

import (
	"github.com/backendengineerark/clients-api/internal/entity"
	"github.com/backendengineerark/clients-api/pkg/customerrors"
	"github.com/backendengineerark/clients-api/pkg/dates"
)

type ClientInputDTO struct {
	Name      string `json:"name"`
	Document  string `json:"document"`
	BirthDate string `json:"birth_date"`
}

type AccountInputDTO struct {
	AccountType    string         `json:"account_type"`
	ClientInputDTO ClientInputDTO `json:"client"`
}

type ClientOutputDTO struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Document  string `json:"document"`
	BirthDate string `json:"birth_date"`
}

type AccountOutputDTO struct {
	Number          string           `json:"number"`
	AccountType     string           `json:"account_type"`
	ClientOutputDTO *ClientOutputDTO `json:"client"`
}

type CreateAccountUseCase struct {
}

func NewCreateAccountUseCase() *CreateAccountUseCase {
	return &CreateAccountUseCase{}
}

func (ca *CreateAccountUseCase) Execute(input AccountInputDTO) (*AccountOutputDTO, []customerrors.Error) {
	client, errors := entity.NewClient(input.ClientInputDTO.Name, input.ClientInputDTO.Document, input.ClientInputDTO.BirthDate)
	if len(errors) > 0 {
		return nil, errors
	}

	account, errors := entity.NewAccount(input.AccountType, *client)
	if len(errors) > 0 {
		return nil, errors
	}

	return &AccountOutputDTO{
		Number:      account.Number,
		AccountType: string(account.AccountType),
		ClientOutputDTO: &ClientOutputDTO{
			Id:        account.Client.Id,
			Name:      account.Client.Name,
			Document:  account.Client.Document,
			BirthDate: dates.DateToString(account.Client.BirthDate),
		},
	}, nil
}
