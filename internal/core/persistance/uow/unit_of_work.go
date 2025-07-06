package uow

import (
	"context"
	"database/sql"
	"fmt"
)

type UnitOfWork interface {
	Do(ctx context.Context, fn func() error) error
	GetExecutor() DBExecutor
}

type unitOfWork struct {
	db *sql.DB
	tx *sql.Tx
}

func NewUnitOfWork(db *sql.DB) UnitOfWork {
	return &unitOfWork{db: db}
}

func (u *unitOfWork) GetExecutor() DBExecutor {
	if u.tx != nil {
		return u.tx
	}
	return u.db
}

func (u *unitOfWork) Do(ctx context.Context, fn func() error) (err error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	u.tx = tx

	defer func() {
		if p := recover(); p != nil {
			_ = u.tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = u.tx.Rollback()
		} else {
			err = u.tx.Commit()
		}
		u.tx = nil
	}()

	err = fn()
	return err
}
