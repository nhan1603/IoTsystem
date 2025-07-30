package test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func LoadSqlTestFile(t *testing.T, db *sql.Tx, sqlfile string) {
	b, err := os.ReadFile(sqlfile)
	require.NoError(t, err)

	_, err = db.Exec(string(b))
	require.NoError(t, err)
}
