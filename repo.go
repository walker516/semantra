package repo

import (
	"api/internal/core"
	"api/internal/repo/postgres"

	"github.com/jmoiron/sqlx"
)

type Repositories struct {
	Feature core.FeatureRepository
	// Tag core.TagRepository
}

func NewRepositories(pg *sqlx.DB) *Repositories {
	return &Repositories{
		Feature: postgres.NewFeatureRepository(pg),
	}
}
