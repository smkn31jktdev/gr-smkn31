package util

import (
	"fmt"
	"math"
)

const earthRadius = 6371000.0

// Haversine menghitung jarak dua titik koordinat dalam meter.
func Haversine(lat1, lng1, lat2, lng2 float64) float64 {
	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadius * c
}

// ValidateRadius memvalidasi jarak siswa dari sekolah.
func ValidateRadius(studentLat, studentLng, sekolahLat, sekolahLng, maxMeter float64) (float64, error) {
	dist := Haversine(studentLat, studentLng, sekolahLat, sekolahLng)
	if dist > maxMeter {
		return dist, fmt.Errorf("jarak %.1f meter melebihi batas %.0f meter", dist, maxMeter)
	}
	return dist, nil
}
