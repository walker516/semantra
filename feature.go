package core

import (
	"api/pkg/nullutil"
	"time"
)

type Feature struct {
	ID        string            `db:"st02_feature_id"  json:"id"`
	ModelID   string            `db:"st01_model_id"    json:"model_id"`
	Geom      string            `db:"geom"`
	Props     nullutil.JSONBMap `db:"st02_props"       json:"props"`
	ObsTime   *time.Time        `db:"st02_obs_time"    json:"obs_time,omitempty"`
	CreatedAt time.Time         `db:"st02_created_at"  json:"created_at"`
	UpdatedAt time.Time         `db:"st02_updated_at"  json:"updated_at"`
}
