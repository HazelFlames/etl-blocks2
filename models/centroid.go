package models

import (
	geojson "github.com/paulmach/go.geojson"
	"github.com/spatial-go/geoos/encoding/wkt"
	"github.com/spatial-go/geoos/space"
)

// Centroid for only 4 vertices
func Centroid(polygon [][][]float64) *geojson.Geometry {
	minLon, minLat := polygon[0][0][0], polygon[0][0][1]
	maxLon, maxLat := polygon[0][3][0], polygon[0][3][1]

	centroidLon := (maxLon-minLon)/2.0 + minLon
	centroidLat := (maxLat-minLat)/2.0 + minLat

	return geojson.NewPointGeometry([]float64{centroidLon, centroidLat})
}

func CentroidText(centroid *geojson.Geometry) string {
	point := space.Point{centroid.Point[0], centroid.Point[1]}
	return wkt.MarshalString(point)
}
