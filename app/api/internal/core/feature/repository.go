package feature

type FeatureRepository interface {
	FindByID(id FeatureID) (*Feature, error)
	Save(f *Feature) error
	ListByBBox(bbox [4]float64) ([]*Feature, error)
}