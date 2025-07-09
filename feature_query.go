package query

import (
	"api/internal/core"
	"context"
)

type FeatureQueryFlow interface {
	GetByID(ctx context.Context, id string) (*core.Feature, error)
}

type featureQueryFlow struct{ repo core.FeatureRepository }

func NewFeatureQuery(r core.FeatureRepository) FeatureQueryFlow {
	return &featureQueryFlow{repo: r}
}

func (u *featureQueryFlow) GetByID(ctx context.Context, id string) (*core.Feature, error) {
	// validation UUID 等はここで
	return u.repo.FindByID(ctx, id)
}
