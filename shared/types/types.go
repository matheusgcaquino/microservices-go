package types

type Routes struct {
	Routes []*Route
}

type Route struct {
	Distance float64
	Duration float64
	Geometry *Geometry
}

type Geometry struct {
	Coordinates []*Coordinates
}

type Coordinates struct {
	Longitude float64
	Latitude  float64
}
