package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/backendengineerark/clients-api/internal/entity"
	"github.com/backendengineerark/clients-api/pkg/conversions"
	"github.com/backendengineerark/clients-api/pkg/customerrors"
	"github.com/backendengineerark/clients-api/pkg/customlogs"
	"github.com/backendengineerark/clients-api/pkg/dates"
	"github.com/backendengineerark/clients-api/pkg/events"
	"github.com/backendengineerark/clients-api/pkg/uow"
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
	AccountStatus   string           `json:"account_status"`
	ClientOutputDTO *ClientOutputDTO `json:"client"`
	CreatedAt       time.Time        `json:"created_at"`
}

type CreateAccountUseCase struct {
	Uow                 uow.Uow
	ClientRepository    entity.ClientRepositoryInterface
	AccountRepository   entity.AccountRepositoryInterface
	AccountCreatedEvent events.EventInterface
	EventDispatcher     events.EventDispatcherInterface
}

func NewCreateAccountUseCase(uow *uow.Uow, clientRepository entity.ClientRepositoryInterface, accountRepository entity.AccountRepositoryInterface, accountCreatedEvent events.EventInterface, eventDispatcher events.EventDispatcherInterface) *CreateAccountUseCase {
	return &CreateAccountUseCase{
		Uow:                 *uow,
		ClientRepository:    clientRepository,
		AccountRepository:   accountRepository,
		AccountCreatedEvent: accountCreatedEvent,
		EventDispatcher:     eventDispatcher,
	}
}

func (ca *CreateAccountUseCase) Execute(ctx context.Context, input AccountInputDTO) (*AccountOutputDTO, []customerrors.Error, error) {
	logger := customlogs.GetContextLogger(ctx)
	logger.Printf("Try to validate client")

	client, errors := entity.NewClient(input.ClientInputDTO.Name, input.ClientInputDTO.Document, input.ClientInputDTO.BirthDate)
	if len(errors) > 0 {
		logger.Printf("Fail to instance client because %s", conversions.StructToJsonIgnoreErrors(ctx, errors))
		return nil, errors, nil
	}
	logger.Printf("Success to instance client")

	logger.Printf("Try to validate account")
	account, errors := entity.NewAccount(input.AccountType, *client)
	if len(errors) > 0 {
		logger.Printf("Fail to instance account because %s", conversions.StructToJsonIgnoreErrors(ctx, errors))
		return nil, errors, nil
	}
	logger.Printf("Success to instance account with number %s", account.Number)

	logger.Printf("Validate if client already exists by document")
	clientExists, err := ca.ClientRepository.ExistsByDocument(input.ClientInputDTO.Document)
	if err != nil {
		logger.Printf("Fail to get client by document %s", err)
		return nil, nil, err
	}
	if clientExists {
		logger.Printf("Client already exists with that document")
		return nil, []customerrors.Error{*customerrors.NewError(customerrors.CLIENT_ALREADY_EXISTS, fmt.Sprintf("Client already exists with document %s", input.ClientInputDTO.Document))}, nil
	}
	logger.Printf("Client not exists with that document")

	err = ca.Persist(ctx, client, account)
	if err != nil {
		logger.Printf("Fail to persist %s", err)
		return nil, nil, err
	}

	logger.Printf("Success to save client and account")
	accountOutputDto := &AccountOutputDTO{
		Number:        account.Number,
		AccountType:   string(account.AccountType),
		AccountStatus: string(account.AccountStatus),
		CreatedAt:     account.CreatedAt,
		ClientOutputDTO: &ClientOutputDTO{
			Id:        account.Client.Id,
			Name:      account.Client.Name,
			Document:  account.Client.Document,
			BirthDate: dates.DateToString(account.Client.BirthDate),
			CreatedAt: account.Client.CreatedAt,
		},
	}

	ca.AccountCreatedEvent.SetPayload(accountOutputDto)
	ca.EventDispatcher.Dispatch(ctx, ca.AccountCreatedEvent)

	return accountOutputDto, nil, nil
}

func (ca CreateAccountUseCase) Persist(ctx context.Context, client *entity.Client, account *entity.Account) error {
	logger := customlogs.GetContextLogger(ctx)

	return ca.Uow.Execute(ctx, func(tx *sql.Tx) error {
		logger.Printf("Try to save a client")
		if err := ca.ClientRepository.Save(tx, client); err != nil {
			return err
		}
		logger.Printf("Success to save client")

		logger.Printf("Try to save account")
		if err := ca.AccountRepository.Save(tx, account); err != nil {
			return err
		}
		logger.Printf("Success to save account")
		return nil
	})
}
