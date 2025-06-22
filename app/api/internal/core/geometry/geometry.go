package geometry

type GeometryType string

const (
	PointType   GeometryType = "Point"
	LineType    GeometryType = "Line"
	PolygonType GeometryType = "Polygon"
)

type Geometry interface {
	Type() GeometryType
	ToGeoJSON() map[string]any
}