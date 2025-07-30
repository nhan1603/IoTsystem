package test

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/nhan1603/IoTsystem/api/internal/appconfig/db/pg"
	"github.com/stretchr/testify/require"
)

// appDB caches a pg connection for reuse
var appDB *sql.DB

// WithTxDB provides callback with a `sql transaction` for running pg related tests
func WithTxDB(t *testing.T, callback func(*sql.Tx)) {
	if appDB == nil {
		var err error
		appDB, err = pg.Connect(os.Getenv("PG_URL"))

		require.NoError(t, err)
	}

	tx, err := appDB.BeginTx(context.Background(), nil)
	require.NoError(t, err)

	// rollback transaction after finish test
	defer tx.Rollback()

	callback(tx)
}
