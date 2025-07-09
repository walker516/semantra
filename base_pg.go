package postgres

import (
	"api/pkg/logutil"
	"api/pkg/tmplutil"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// PG is injected to concrete repos.  thread-safe.
type PG struct {
	DB       *sqlx.DB
	Template *tmplutil.SQLTemplateEngine
}

func NewPG(db *sqlx.DB, tmplPath string) *PG {
	return &PG{
		DB:       db,
		Template: tmplutil.NewSQLTemplateEngine(tmplPath),
	}
}

func (pg *PG) query(ctx context.Context, section string, params map[string]any, scan func(*sqlx.Rows) error) error {
	ctx, span := otel.Tracer("repo").Start(ctx, section)
	defer span.End()

	sqlTmpl, _, err := pg.Template.RenderQuery(section, params)
	if err != nil {
		return err
	}

	sqlNamed, args, err := sqlx.Named(sqlTmpl, params)
	if err != nil {
		return err
	}

	sqlRebound := sqlx.Rebind(sqlx.DOLLAR, sqlNamed)
	pg.logSQL(sqlRebound, args)

	rows, err := pg.DB.QueryxContext(ctx, sqlRebound, args...)
	if err != nil {
		fmt.Printf("[query] DB QueryxContext ERROR: %+v\n", err)
		return err
	}
	defer rows.Close()

	scanErr := scan(rows)
	fmt.Printf("[query] Scan func result: %+v\n", scanErr)
	return scanErr
}

func (pg *PG) exec(ctx context.Context, section string, params map[string]any) (int64, error) {
	ctx, span := otel.Tracer("repo").Start(ctx, section)
	defer span.End()

	sqlTmpl, _, err := pg.Template.RenderQuery(section, params)
	if err != nil {
		return 0, err
	}

	sqlNamed, args, err := sqlx.Named(sqlTmpl, params)
	if err != nil {
		return 0, err
	}

	sqlRebound := sqlx.Rebind(sqlx.DOLLAR, sqlNamed)
	pg.logSQL(sqlRebound, args)

	res, err := pg.DB.ExecContext(ctx, sqlRebound, args...)
	if err != nil {
		return 0, err
	}

	affected, _ := res.RowsAffected()
	span.SetAttributes(attribute.Int64("rows_affected", affected))
	return affected, nil
}

func (pg *PG) logSQL(q string, args any) {
	logutil.Info("SQL executed",
		zap.String("query", q),
		zap.Any("args", args),
	)
}
