package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/nhan1603/IoTsystem/api/internal/repository/iotsystem"
	"github.com/nhan1603/IoTsystem/api/internal/repository/user"
	pkgerrors "github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type Registry interface {
	User() user.Repository
	IoT() iotsystem.Repository
	DoInTx(ctx context.Context, txFunc TxFunc) error
}

// New returns an implementation instance which satisfying Registry
func New(pgConn *sql.DB) Registry {
	return impl{
		user:   user.New(pgConn),
		iot:    iotsystem.New(pgConn),
		pgConn: pgConn,
	}
}

type impl struct {
	user   user.Repository
	iot    iotsystem.Repository
	txExec boil.Transactor
	pgConn *sql.DB
}

// TxFunc is a function that can be executed in a transaction
type TxFunc func(txRegistry Registry) error

// User returns user repo
func (i impl) User() user.Repository {
	return i.user
}

func (i impl) IoT() iotsystem.Repository {
	return i.iot
}

// DoInTx handles db operations in a transaction
func (i impl) DoInTx(ctx context.Context, txFunc TxFunc) error {
	if i.txExec != nil {
		return errors.New("db tx nested in db tx")
	}

	tx, err := i.pgConn.BeginTx(ctx, nil)
	if err != nil {
		return pkgerrors.WithStack(err)
	}

	var committed bool
	defer func() {
		if committed {
			return
		}

		_ = tx.Rollback()
	}()

	newI := impl{
		user:   user.New(tx),
		iot:    iotsystem.New(tx),
		txExec: tx,
	}

	if err = txFunc(newI); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return pkgerrors.WithStack(err)
	}

	committed = true

	return nil
}
