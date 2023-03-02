package event

import "time"

type AccountCreated struct {
	Name    string
	Payload interface{}
}

func NewAccountCreated() *AccountCreated {
	return &AccountCreated{
		Name: "AccountCreated",
	}
}

func (ac *AccountCreated) GetName() string {
	return ac.Name
}

func (ac *AccountCreated) GetDateTime() time.Time {
	return time.Now()
}

func (ac *AccountCreated) GetPayload() interface{} {
	return ac.Payload
}

func (ac *AccountCreated) SetPayload(payload interface{}) {
	ac.Payload = payload
}
