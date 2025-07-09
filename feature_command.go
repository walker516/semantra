package command

import (
	"api/internal/core"
	"context"

	"github.com/google/uuid"
)

type FeatureCommandFlow interface {
	Create(ctx context.Context, f *core.Feature) (string, error)
	Update(ctx context.Context, f *core.Feature) error
	Delete(ctx context.Context, id string) error
}

type featureCommandFlow struct{ repo core.FeatureRepository }

func NewFeatureCommand(r core.FeatureRepository) FeatureCommandFlow {
	return &featureCommandFlow{repo: r}
}

func (u *featureCommandFlow) Create(ctx context.Context, f *core.Feature) (string, error) {
	// 生成側でUUIDv4 を明示生成するパターン
	if f.ID == "" {
		f.ID = uuid.New().String()
	}
	return f.ID, u.repo.Create(ctx, f)
}

func (u *featureCommandFlow) Update(ctx context.Context, f *core.Feature) error {
	return u.repo.Update(ctx, f)
}

func (u *featureCommandFlow) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}
