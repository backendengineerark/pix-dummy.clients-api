package entity

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/backendengineerark/clients-api/pkg/customerrors"
)

type AccountType string

const (
	ContaCorrente AccountType = "CONTA_CORRENTE"
	ContaPoupanca AccountType = "CONTA_POUPANCA"
)

type AccountStatus string

const (
	AccountActive   AccountStatus = "OPEN"
	AccountCanceled AccountStatus = "CANCELED"
)

type Account struct {
	Number        string
	AccountType   AccountType
	AccountStatus AccountStatus
	Client        Client
	CreatedAt     time.Time
}

func NewAccount(accountType string, client Client) (*Account, []customerrors.Error) {
	errors := []customerrors.Error{}

	accountTypeCreated, err := CreateAccountType(accountType)
	if err != nil {
		return nil, append(errors, *err)
	}

	account := &Account{
		Number:        GenerateAccountNumber(),
		AccountType:   accountTypeCreated,
		AccountStatus: AccountActive,
		Client:        client,
		CreatedAt:     time.Now(),
	}

	validationErrors := account.IsValid()
	if len(validationErrors) > 0 {
		errors = append(errors, validationErrors...)
		return nil, errors
	}
	return account, nil
}

func (acc *Account) IsValid() []customerrors.Error {
	errorsList := []customerrors.Error{}

	if acc.Number == "" {
		errorsList = append(errorsList, *customerrors.NewError(customerrors.INVALID_PARAM, "Number is required"))
	}

	if acc.AccountType == "" {
		errorsList = append(errorsList, *customerrors.NewError(customerrors.INVALID_PARAM, "Account type is required"))
	}

	if acc.AccountStatus == "" {
		errorsList = append(errorsList, *customerrors.NewError(customerrors.INVALID_PARAM, "Account status is required"))
	}

	if acc.Client == (Client{}) {
		errorsList = append(errorsList, *customerrors.NewError(customerrors.INVALID_PARAM, "Client is required"))
	}

	if acc.CreatedAt.IsZero() {
		errorsList = append(errorsList, *customerrors.NewError(customerrors.INVALID_PARAM, "Created at is required"))
	}

	return errorsList
}

func CreateAccountType(accountType string) (AccountType, *customerrors.Error) {

	if accountType == string(ContaCorrente) {
		return ContaCorrente, nil
	}

	if accountType == string(ContaPoupanca) {
		return ContaPoupanca, nil
	}

	return ContaCorrente, customerrors.NewError(customerrors.INVALID_PARAM, fmt.Sprintf("Account type should be %s or %s", string(ContaCorrente), string(ContaPoupanca)))
}

func GenerateAccountNumber() string {
	min := 10000000
	max := 99999999
	return fmt.Sprint((rand.Intn(max-min) + min))
}
