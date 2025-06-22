package flow

import (
	"api/internal/core/feature"
)

type FeatureUsecase struct{ repo feature.FeatureRepository }

func NewFeatureUsecase(r feature.FeatureRepository) *FeatureUsecase { return &FeatureUsecase{repo: r} }

func (uc *FeatureUsecase) Get(id feature.FeatureID) (*feature.Feature, error) {
	return uc.repo.FindByID(id)
}
