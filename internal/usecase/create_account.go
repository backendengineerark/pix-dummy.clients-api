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
	Db                  sql.DB
	ClientRepository    entity.ClientRepositoryInterface
	AccountRepository   entity.AccountRepositoryInterface
	AccountCreatedEvent events.EventInterface
	EventDispatcher     events.EventDispatcherInterface
}

func NewCreateAccountUseCase(db sql.DB, clientRepository entity.ClientRepositoryInterface, accountRepository entity.AccountRepositoryInterface, accountCreatedEvent events.EventInterface, eventDispatcher events.EventDispatcherInterface) *CreateAccountUseCase {
	return &CreateAccountUseCase{
		Db:                  db,
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
		logger.Printf("Fail to validate client because %s", conversions.StructToJsonIgnoreErrors(ctx, errors))
		return nil, errors, nil
	}
	logger.Printf("Success to validate client")

	logger.Printf("Try to validate account")
	account, errors := entity.NewAccount(input.AccountType, *client)
	if len(errors) > 0 {
		logger.Printf("Fail to validate account because %s", conversions.StructToJsonIgnoreErrors(ctx, errors))
		return nil, errors, nil
	}
	logger.Printf("Success to validate account")

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

	logger.Printf("Try to start a transaction")
	tx, err := ca.Db.BeginTx(context.Background(), nil)
	if err != nil {
		logger.Printf("Fail to start a transaction because %s", err)
		return err
	}
	logger.Printf("Transaction started")

	logger.Printf("Try to save a client")
	if err := ca.ClientRepository.Save(tx, client); err != nil {
		logger.Printf("Error to save client %s", err)
		tx.Rollback()
		return err
	}
	logger.Printf("Success to save client")

	logger.Printf("Try to save account")
	if err := ca.AccountRepository.Save(tx, account); err != nil {
		logger.Printf("Error to save account, rollback started %s", err)
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
