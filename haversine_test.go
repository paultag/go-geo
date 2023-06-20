package geo_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"pault.ag/go/geo"
)

var testsCases = []struct {
	from           geo.LLA
	to             geo.LLA
	expectedMeters float64
}{
	{
		geo.LLA{Latitude: 22.55, Longitude: 43.12},
		geo.LLA{Latitude: 13.45, Longitude: 100.28},
		6094544.408786774,
	},
	{
		geo.LLA{Latitude: 51.510357, Longitude: -0.116773},
		geo.LLA{Latitude: 38.889931, Longitude: -77.009003},
		5897658.288856054,
	},
}

func TestDistance(t *testing.T) {
	for _, input := range testsCases {
		meters, err := geo.HaversineDistance(input.from, input.to)
		assert.NoError(t, err)
		assert.InEpsilon(t, input.expectedMeters, meters.F64(), 1e-6)
	}
}

func TestDistanceWithAlt(t *testing.T) {
	a := geo.LLA{Latitude: 51.510357, Longitude: -0.116773}
	b := geo.LLA{Latitude: 51.510357, Longitude: -0.116773, Altitude: 10}

	_, err := geo.HaversineDistance(a, b)
	assert.Error(t, err)

	_, err = geo.HaversineDistance(b, a)
	assert.Error(t, err)

	_, err = geo.HaversineDistance(b, b)
	assert.Error(t, err)

	_, err = geo.HaversineDistance(a, a)
	assert.NoError(t, err)
}

func BenchmarkDistance(b *testing.B) {

	from := geo.LLA{Latitude: 22.55, Longitude: 43.12}
	to := geo.LLA{Latitude: 13.45, Longitude: 100.28}
	for i := 0; i < b.N; i++ {
		_, _ = geo.HaversineDistance(from, to)
	}
}
