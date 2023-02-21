package usecase

import (
	"fmt"
	"log"
	"time"

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
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Document  string    `json:"document"`
	BirthDate string    `json:"birth_date"`
	CreatedAt time.Time `json:"created_at"`
}

type AccountOutputDTO struct {
	Number          string           `json:"number"`
	AccountType     string           `json:"account_type"`
	ClientOutputDTO *ClientOutputDTO `json:"client"`
	CreatedAt       time.Time        `json:"created_at"`
}

type CreateAccountUseCase struct {
	ClientRepository  entity.ClientRepositoryInterface
	AccountRepository entity.AccountRepositoryInterface
}

func NewCreateAccountUseCase(clientRepository entity.ClientRepositoryInterface, accountRepository entity.AccountRepositoryInterface) *CreateAccountUseCase {
	return &CreateAccountUseCase{
		ClientRepository:  clientRepository,
		AccountRepository: accountRepository,
	}
}

func (ca *CreateAccountUseCase) Execute(input AccountInputDTO) (*AccountOutputDTO, []customerrors.Error, error) {

	client, errors := entity.NewClient(input.ClientInputDTO.Name, input.ClientInputDTO.Document, input.ClientInputDTO.BirthDate)
	if len(errors) > 0 {
		return nil, errors, nil
	}

	account, errors := entity.NewAccount(input.AccountType, *client)
	if len(errors) > 0 {
		return nil, errors, nil
	}

	clientExists, err := ca.ClientRepository.ExistsByDocument(input.ClientInputDTO.Document)
	if err != nil {
		log.Printf("Fail to get client by document %s", err)
		return nil, nil, err
	}

	if clientExists {
		return nil, []customerrors.Error{*customerrors.NewError(customerrors.CLIENT_ALREADY_EXISTS, fmt.Sprintf("Client already exists with document %s", input.ClientInputDTO.Document))}, nil
	}

	err = ca.Persist(client, account)
	if err != nil {
		return nil, nil, err
	}

	return &AccountOutputDTO{
		Number:      account.Number,
		AccountType: string(account.AccountType),
		CreatedAt:   account.CreatedAt,
		ClientOutputDTO: &ClientOutputDTO{
			Id:        account.Client.Id,
			Name:      account.Client.Name,
			Document:  account.Client.Document,
			BirthDate: dates.DateToString(account.Client.BirthDate),
			CreatedAt: account.Client.CreatedAt,
		},
	}, nil, nil
}

func (ca CreateAccountUseCase) Persist(
	client *entity.Client,
	account *entity.Account) error {

	if err := ca.ClientRepository.Save(client); err != nil {
		fmt.Printf("Error to persist client %s", err)
		return err
	}

	// account.Client.Id = "1"
	if err := ca.AccountRepository.Save(account); err != nil {
		fmt.Printf("Error to persist account %s", err)
		return err
	}

	return nil
}
