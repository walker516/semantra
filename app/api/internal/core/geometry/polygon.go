package geometry

type Polygon struct{ Coords [][][3]float64 }
func (p Polygon) Type() GeometryType { return PolygonType }
func (p Polygon) ToGeoJSON() map[string]any { return map[string]any{"type": "Polygon", "coordinates": p.Coords} }
