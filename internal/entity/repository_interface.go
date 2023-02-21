package entity

type ClientRepositoryInterface interface {
	Save(client *Client) error
	ExistsByDocument(document string) (bool, error)
	// GetAll() (*[]Client, error)
}

type AccountRepositoryInterface interface {
	Save(account *Account) error
	ExistsByNumber(accountNumber string) (bool, error)
	// GetAll() (*[]Account, error)
	// GetById(accountId string) (*Account, error)
}
