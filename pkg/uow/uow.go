package uow

import (
	"context"
	"database/sql"
	"fmt"
)

type RepositoryFactory func(tx *sql.Tx) interface{}

type UowInterface interface {
	Execute(ctx context.Context, fn func(uow *Uow) error) error
}

type Uow struct {
	Db *sql.DB
}

func NewUow(db *sql.DB) *Uow {
	return &Uow{
		Db: db,
	}
}

func (u *Uow) Execute(ctx context.Context, fn func(tx *sql.Tx) error) error {
	conn, err := u.Db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("fail to get a connection %s", err)
	}
	defer conn.Close()

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("fail to start a transaction %s", err)
	}

	err = fn(tx)
	if err != nil {
		errRb := u.rollback(tx)
		if errRb != nil {
			return fmt.Errorf("original error: %s, rollback error: %s", err.Error(), errRb.Error())
		}
		return err
	}
	return u.commitOrRollback(tx)
}

func (u *Uow) rollback(tx *sql.Tx) error {
	err := tx.Rollback()
	if err != nil {
		return fmt.Errorf("fail to rollback transaction %s", err)
	}
	return nil
}

func (u *Uow) commitOrRollback(tx *sql.Tx) error {
	err := tx.Commit()
	if err != nil {
		errRb := u.rollback(tx)
		if errRb != nil {
			return fmt.Errorf("original error: %s, rollback error: %s", err.Error(), errRb.Error())
		}
		return err
	}
	return nil
}
