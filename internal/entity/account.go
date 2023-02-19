package entity

import (
	"errors"
	"fmt"
	"math/rand"
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
}

func NewAccount(accountType string, client Client) (*Account, error) {
	accountTypeCreated, err := CreateAccountType(accountType)
	if err != nil {
		return nil, err
	}
	return &Account{
		Number:      GenerateAccountNumber(),
		AccountType: accountTypeCreated,
		Client:      client,
	}, nil
}

func CreateAccountType(accountType string) (AccountType, error) {

	if accountType == string(ContaCorrente) {
		return ContaCorrente, nil
	}

	if accountType == string(ContaPoupanca) {
		return ContaPoupanca, nil
	}

	return ContaCorrente, errors.New("account type should be " + string(ContaCorrente) + " or " + string(ContaPoupanca))
}

func GenerateAccountNumber() string {
	min := 10000000
	max := 99999999
	return fmt.Sprint((rand.Intn(max-min) + min))
}
