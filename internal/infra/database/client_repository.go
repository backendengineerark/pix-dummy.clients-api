package database

import (
	"database/sql"
	"log"

	"github.com/backendengineerark/clients-api/internal/entity"
)

type ClientRepository struct {
	Db *sql.DB
}

func NewClientRepository(db *sql.DB) *ClientRepository {
	return &ClientRepository{
		Db: db,
	}
}

func (cr *ClientRepository) Save(tx *sql.Tx, client *entity.Client) error {
	stmt, err := tx.Prepare("INSERT INTO clients(id, name, document, birth_date, created_at) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(client.Id, client.Name, client.Document, client.BirthDate, client.CreatedAt)
	if err != nil {
		log.Printf("Fail to save client because %s", err)
		return err
	}
	return nil
}

func (cr *ClientRepository) ExistsByDocument(document string) (bool, error) {
	var exists bool
	if err := cr.Db.QueryRow("SELECT COUNT(*) FROM clients WHERE document = ?", document).Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return exists, nil
}
