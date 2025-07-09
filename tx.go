package dbutil

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// WithTx creates a Tx and passes *sqlx.Tx to fn(); commit or rollback automatically.
func WithTx(ctx context.Context, db *sqlx.DB, fn func(*sqlx.Tx) error) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
