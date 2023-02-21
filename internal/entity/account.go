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

type Account struct {
	Number      string
	AccountType AccountType
	Client      Client
	CreatedAt   time.Time
}

func NewAccount(accountType string, client Client) (*Account, []customerrors.Error) {
	errors := []customerrors.Error{}

	accountTypeCreated, err := CreateAccountType(accountType)
	if err != nil {
		return nil, append(errors, *err)
	}
	return &Account{
		Number:      GenerateAccountNumber(),
		AccountType: accountTypeCreated,
		Client:      client,
		CreatedAt:   time.Now(),
	}, nil
}

func CreateAccountType(accountType string) (AccountType, *customerrors.Error) {

	if accountType == string(ContaCorrente) {
		return ContaCorrente, nil
	}

	if accountType == string(ContaPoupanca) {
		return ContaPoupanca, nil
	}

	return ContaCorrente, customerrors.NewError(customerrors.INVALID, fmt.Sprintf("Account type should be %s or %s", string(ContaCorrente), string(ContaPoupanca)))
}

func GenerateAccountNumber() string {
	min := 10000000
	max := 99999999
	return fmt.Sprint((rand.Intn(max-min) + min))
}
