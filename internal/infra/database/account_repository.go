package database

import (
	"database/sql"

	"github.com/backendengineerark/clients-api/internal/entity"
)

type AccountRepository struct {
	Db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{
		Db: db,
	}
}

func (ar *AccountRepository) Save(account *entity.Account) error {
	stmt, err := ar.Db.Prepare("INSERT INTO accounts(number, account_type, client_id, created_at) VALUES(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(account.Number, string(account.AccountType), account.Client.Id, account.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (cr *AccountRepository) ExistsByNumber(accountNumber string) (bool, error) {
	var exists bool
	if err := cr.Db.QueryRow("SELECT COUNT(*) FROM accounts WHERE number = ?", accountNumber).Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return exists, nil
}
