package model

type DriverLocation struct {
	DriverID string   `bson:"driverId" json:"driverId"`
	Location GeoPoint `bson:"location" json:"location"`
}

type GeoPoint struct {
	Type        string     `bson:"type" json:"type"`               // Point
	Coordinates [2]float64 `bson:"coordinates" json:"coordinates"` // [latitude, longitude] but MongoDB accepts [longitude, latitude] if sending individually
}

type DriverWithDistance struct {
	DriverLocation `bson:",inline" json:",inline"`
	DistanceMeters float64 `bson:"distanceMeters" json:"distanceMeters"`
}
