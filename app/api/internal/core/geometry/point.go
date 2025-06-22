package geometry

type Point struct{ Lng, Lat, Alt float64 }
func (p Point) Type() GeometryType { return PointType }
func (p Point) ToGeoJSON() map[string]any {
	coords := []float64{p.Lng, p.Lat}
	if p.Alt != 0 { coords = append(coords, p.Alt) }
	return map[string]any{"type": "Point", "coordinates": coords}
}
