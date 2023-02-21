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

func (ar *AccountRepository) Save(tx *sql.Tx, account *entity.Account) error {
	stmt, err := tx.Prepare("INSERT INTO accounts(number, account_type, account_status, client_id, created_at) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(account.Number, account.AccountType, account.AccountStatus, account.Client.Id, account.CreatedAt)
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
