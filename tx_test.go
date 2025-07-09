package dbutil_test

import (
	"context"
	"testing"

	"api/pkg/dbutil"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestWithTx_Commit(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	mock.ExpectBegin()
	mock.ExpectCommit()

	err = dbutil.WithTx(context.Background(), sqlxDB, func(tx *sqlx.Tx) error {
		return nil // no error → should commit
	})
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestWithTx_Rollback(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	mock.ExpectBegin()
	mock.ExpectRollback()

	err = dbutil.WithTx(context.Background(), sqlxDB, func(tx *sqlx.Tx) error {
		return assertAnError()
	})
	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func assertAnError() error {
	return sqlmock.ErrCancelled
}
