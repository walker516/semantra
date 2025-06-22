package feature

import "api/internal/core/geometry"

type FeatureKind string

const (
	TileKind     FeatureKind = "tile"
	RoadKind     FeatureKind = "road"
	BuildingKind FeatureKind = "building"
)

type FeatureID string

type Feature struct {
	ID         FeatureID         `json:"id"`
	Kind       FeatureKind       `json:"kind"`
	Geometry   geometry.Geometry `json:"geometry"`
	Properties map[string]any    `json:"properties"`
}