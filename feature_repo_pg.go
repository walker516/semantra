package postgres

import (
	"api/internal/core"
	"api/pkg/errmsg"
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/jmoiron/sqlx"
)

type pgRepo struct{ *PG }

func NewFeatureRepository(db *sqlx.DB) core.FeatureRepository {
	return &pgRepo{
		PG: NewPG(db, filepath.Join("templates", "sql", "feature", "queries.sql")),
	}
}

func (r *pgRepo) FindByID(ctx context.Context, id string) (*core.Feature, error) {
	params := map[string]any{"FeatureID": id}
	var feat core.Feature

	err := r.query(ctx, "FindByID", params, func(rows *sqlx.Rows) error {
		if rows.Next() {
			if err := rows.StructScan(&feat); err != nil {
				fmt.Printf("[repo] SCAN ERROR: %+v\n", err)
				return errmsg.Wrap("ERR_INTERNAL_SERVER", "Failed to scan user data", err)
			}
			fmt.Println("[repo] FOUND FEATURE:", feat)
			return nil
		}

		fmt.Println("[repo] NO ROWS found, returning ERR_NOT_FOUND")
		return errmsg.New("ERR_NOT_FOUND", fmt.Sprintf("Feature %s not found", id))
	})

	fmt.Printf("[repo] Final returned error: %#v\n", err)

	if err != nil {
		return nil, err
	}
	return &feat, nil
}

func (r *pgRepo) Create(ctx context.Context, f *core.Feature) error {
	params := map[string]any{
		"FeatureID": f.ID,
		"ModelID":   f.ModelID,
		"GeomJSON":  f.Geom,  // 事前にGeoJSONに整形済み
		"Props":     f.Props, // map[string]interface{}
		"ObsTime":   f.ObsTime,
	}
	_, err := r.exec(ctx, "CreateFeature", params)
	return err
}

func (r *pgRepo) Update(ctx context.Context, f *core.Feature) error {
	params := map[string]any{
		"FeatureID": f.ID,
		"Props":     f.Props, // 差分のみを JSONB で更新も可
		"UpdatedAt": time.Now(),
	}
	affected, err := r.exec(ctx, "UpdateFeature", params)
	if err != nil {
		return err
	}
	if affected == 0 {
		return errmsg.New("ERR_NOT_FOUND", fmt.Sprintf("Feature %s not found", f.ID))
	}
	return nil
}

func (r *pgRepo) Delete(ctx context.Context, id string) error {
	params := map[string]any{"FeatureID": id}
	affected, err := r.exec(ctx, "DeleteFeature", params)
	if err != nil {
		return err
	}
	if affected == 0 {
		return errmsg.New("ERR_NOT_FOUND", fmt.Sprintf("Feature %s not found", id))
	}
	return nil
}
