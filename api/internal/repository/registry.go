package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/gocql/gocql"
	"github.com/nhan1603/IoTsystem/api/internal/repository/cassiot"
	"github.com/nhan1603/IoTsystem/api/internal/repository/iotsystem"
	"github.com/nhan1603/IoTsystem/api/internal/repository/user"
	pkgerrors "github.com/pkg/errors"
	"github.com/scylladb/gocqlx/v2"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type Registry interface {
	User() user.Repository
	IoT() iotsystem.Repository
	DoInTx(ctx context.Context, txFunc TxFunc) error
}

// newPostGres returns an implementation instance which satisfying Registry
func newPostGres(pgConn *sql.DB) Registry {
	return impl{
		Backend: BackendPostgres,
		user:    user.New(pgConn),
		iot:     iotsystem.New(pgConn),
		pgConn:  pgConn,
	}
}

// newCassandraImpl builds a Cassandra registry (session-bound).
func newCassandraImpl(sess *gocql.Session) impl {
	return impl{
		Backend:     BackendCassandra,
		cassSession: sess,
		// Instantiate repos bound to the session
		user: nil,
		iot:  cassiot.NewCassandra(gocqlx.NewSession(sess)),
	}
}

type impl struct {
	Backend     Backend
	user        user.Repository
	iot         iotsystem.Repository
	txExec      boil.Transactor
	pgConn      *sql.DB
	cassSession *gocql.Session
	inCassBatch bool
}

// TxFunc is a function that can be executed in a transaction
type TxFunc func(txRegistry Registry) error

// Add new type for Cassandra batch operations
type CassandraBatchFunc func(session *gocqlx.Session) error

// User returns user repo
func (i impl) User() user.Repository {
	return i.user
}

func (i impl) IoT() iotsystem.Repository {
	return i.iot
}

// DoInTx handles db operations in a transaction
func (i impl) DoInTx(ctx context.Context, txFunc TxFunc) error {
	switch i.Backend {
	case BackendPostgres:
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
	case BackendCassandra:
		// prevent nested batch scopes
		if i.inCassBatch {
			return errors.New("cassandra batch nested in cassandra batch")
		}
		if i.cassSession == nil {
			return errors.New("cassandra session not initialized")
		}

		batch := i.cassSession.NewBatch(gocql.UnloggedBatch).WithContext(ctx)

		// child impl carries the "batch in ctx"
		child := impl{
			Backend:     i.Backend,
			pgConn:      i.pgConn,
			txExec:      nil,
			cassSession: i.cassSession,
			inCassBatch: true,

			// Cassandra repos can stay the same instance if they read ctx,
			// or you can re-new them if you keep per-impl state. Example keeps them:
			user: i.user,
			iot:  i.iot,
		}

		if err := txFunc(child); err != nil {
			return err
		}
		if err := i.cassSession.ExecuteBatch(batch); err != nil {
			return fmt.Errorf("execute cassandra batch: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("unknown backend: %s", i.Backend)
	}
}
