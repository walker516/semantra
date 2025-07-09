package core

import "context"

// FeatureRepository は DB / ES など永続層の抽象
type FeatureRepository interface {
	FindByID(ctx context.Context, id string) (*Feature, error)
	Create(ctx context.Context, f *Feature) error
	Update(ctx context.Context, f *Feature) error
	Delete(ctx context.Context, id string) error
	// TODO: Create / Update / Search など後で追加
}
