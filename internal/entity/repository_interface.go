package entity

import "database/sql"

type ClientRepositoryInterface interface {
	Save(tx *sql.Tx, client *Client) error
	ExistsByDocument(document string) (bool, error)
	// GetAll() (*[]Client, error)
}

type AccountRepositoryInterface interface {
	Save(tx *sql.Tx, account *Account) error
	ExistsByNumber(accountNumber string) (bool, error)
	// GetAll() (*[]Account, error)
	// GetById(accountId string) (*Account, error)
}
