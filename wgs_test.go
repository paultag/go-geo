package geo_test

import (
	"testing"

	"pault.ag/go/geo"

	"github.com/stretchr/testify/assert"
)

func TestWGS84LLAtoXYZ(t *testing.T) {
	wgs84 := geo.WGS84()
	ref := geo.LLA{
		Latitude:  34.00000048,
		Longitude: -117.3335693,
		Altitude:  251.702,
	}

	refx := wgs84.LLAToXYZ(ref)

	assert.InEpsilon(t, -2430601.8, refx.X.F64(), 0.1)
	assert.InEpsilon(t, -4702442.7, refx.Y.F64(), 0.1)
	assert.InEpsilon(t, 3546587.4, refx.Z.F64(), 0.1)
}

func TestWGS84ConversionCycle(t *testing.T) {
	wgs84 := geo.WGS84()

	ref := geo.LLA{
		Latitude:  38.897957,
		Longitude: -77.036560,
		Altitude:  30,
	}

	position := geo.LLA{
		Latitude:  38.8709455,
		Longitude: -77.0552551,
		Altitude:  100,
	}

	positionx := wgs84.LLAToXYZ(position)
	positionenu := wgs84.XYZToENU(ref, positionx)

	positionxx1 := wgs84.ENUToXYZ(ref, positionenu)
	position1 := wgs84.XYZToLLA(positionxx1)

	assert.InEpsilon(t, float64(position.Latitude), float64(position1.Latitude), 1e-7)
	assert.InEpsilon(t, float64(position.Longitude), float64(position1.Longitude), 1e-7)
	assert.InEpsilon(t, float64(position.Altitude), float64(position1.Altitude), 1e-7)

}

func TestWGS84LLAtoENU(t *testing.T) {
	wgs84 := geo.WGS84()
	ref := geo.LLA{
		Latitude:  34.00000048,
		Longitude: -117.3335693,
		Altitude:  251.702,
	}

	refx := wgs84.LLAToXYZ(ref)

	point := geo.XYZ{
		X: refx.X + 1,
		Y: refx.Y,
		Z: refx.Z,
	}
	pointenu := wgs84.XYZToENU(ref, point)
	assert.InEpsilon(t, 0.88834836, pointenu.East.F64(), 0.1)
	assert.InEpsilon(t, 0.25676467, pointenu.North.F64(), 0.1)
	assert.InEpsilon(t, -0.38066927, pointenu.Up.F64(), 0.1)

	point = geo.XYZ{
		X: refx.X,
		Y: refx.Y + 1,
		Z: refx.Z,
	}
	pointenu = wgs84.XYZToENU(ref, point)
	assert.InEpsilon(t, -0.45917011, pointenu.East.F64(), 0.1)
	assert.InEpsilon(t, 0.49675810, pointenu.North.F64(), 0.1)
	assert.InEpsilon(t, -0.73647416, pointenu.Up.F64(), 0.1)

	point = geo.XYZ{
		X: refx.X,
		Y: refx.Y,
		Z: refx.Z + 1,
	}
	pointenu = wgs84.XYZToENU(ref, point)
	assert.Equal(t, 0.00000000, pointenu.East.F64())
	assert.InEpsilon(t, 0.82903757, pointenu.North.F64(), 0.1)
	assert.InEpsilon(t, 0.55919291, pointenu.Up.F64(), 0.1)
}

func TestENUAERRoundTrip(t *testing.T) {
	enu := geo.ENU{
		East:  10,
		North: 20,
		Up:    30,
	}

	enu1 := enu.AER().ENU()

	assert.InEpsilon(t, 10, enu1.East.F64(), 1e-6)
	assert.InEpsilon(t, 20, enu1.North.F64(), 1e-6)
	assert.InEpsilon(t, 30, enu1.Up.F64(), 1e-6)

}

func TestWGS84AERENURoundTrip(t *testing.T) {
	ref := geo.LLA{
		Latitude:  38.897957,
		Longitude: -77.036560,
		Altitude:  30,
	}

	position := geo.LLA{
		Latitude:  38.8709455,
		Longitude: -77.0552551,
		Altitude:  30,
	}

	wgs84 := geo.WGS84()
	positionx := wgs84.LLAToXYZ(position)
	positionenu := wgs84.XYZToENU(ref, positionx)

	aed := positionenu.AER()
	positionenu1 := aed.ENU()

	positionx1 := wgs84.ENUToXYZ(ref, positionenu1)
	position1 := wgs84.XYZToLLA(positionx1)

	assert.InEpsilon(t, float64(position.Latitude), float64(position1.Latitude), 1e-7)
	assert.InEpsilon(t, float64(position.Longitude), float64(position1.Longitude), 1e-7)
	assert.InEpsilon(t, float64(position.Altitude), float64(position1.Altitude), 1e-7)
}
