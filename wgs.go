// {{{ Copyright (c) Paul R. Tagliamonte <paultag@gmail.com> 2020-2021
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE. }}}

package geo

import (
	"math"
)

var (
	// wgs84A is the earth semimajor axis in Meters
	wgs84A float64 = 6378137.0

	// wgs84B is the earth semiminor axis in Meters
	wgs84B float64 = 6356752.314245

	// wgs84F is the Ellipsoid "flatness"
	wgs84F float64 = (wgs84A - wgs84B) / wgs84A

	// wgs84FInv is the inverse of wgs84F
	wgs84FInv float64 = (1.0 / wgs84F)

	wgs84ASq float64 = wgs84A * wgs84A
	wgs84BSq float64 = wgs84B * wgs84B
	wgs84ESq float64 = wgs84F * (2 - wgs84F)
)

// WGS84 will return the CoordinateSystem for the WSG84 system of Coordinates.
//
// WGS84 is the US Government Coordinate System maintained by the NGA. This
// is what is used by GPS.
func WGS84() CoordinateSystem {
	return wgs84{}
}

type wgs84 struct{}

func (w wgs84) XYZToLLA(x XYZ) LLA {

	var (
		eps   = wgs84ESq / (1 - wgs84ESq)
		p     = math.Sqrt((x.X*x.X + x.Y*x.Y).F64())
		q     = math.Atan2((x.Z.F64() * wgs84A), (p * wgs84B))
		sinQ  = math.Sin(q)
		cosQ  = math.Cos(q)
		sinQ3 = sinQ * sinQ * sinQ
		cosQ3 = cosQ * cosQ * cosQ

		phi = math.Atan2(
			(x.Z.F64() + eps*wgs84B*sinQ3),
			(p - wgs84ESq*wgs84A*cosQ3),
		)
		lambda = math.Atan2(x.Y.F64(), x.X.F64())
		v      = wgs84A / math.Sqrt(1.0-wgs84ESq*math.Sin(phi)*math.Sin(phi))
		h      = Meters((p / math.Cos(phi)) - v)
	)

	return LLA{
		Latitude:  Radians(phi).Degrees(),
		Longitude: Radians(lambda).Degrees(),
		Altitude:  h,
	}
}

func (w wgs84) LLAToXYZ(l LLA) XYZ {
	var (
		lambda = l.Latitude.Radians().F64()
		phi    = l.Longitude.Radians().F64()

		sinLambda = math.Sin(lambda)
		cosLambda = math.Cos(lambda)
		sinPhi    = math.Sin(phi)
		cosPhi    = math.Cos(phi)

		n = wgs84A / math.Sqrt(1-wgs84ESq*sinLambda*sinLambda)
	)

	return XYZ{
		X: Meters((l.Altitude.F64() + n) * cosLambda * cosPhi),
		Y: Meters((l.Altitude.F64() + n) * cosLambda * sinPhi),
		Z: Meters((l.Altitude.F64() + (1-wgs84ESq)*n) * sinLambda),
	}
}

func (w wgs84) LLAToENU(ref, lla LLA) ENU {
	xyz := w.LLAToXYZ(lla)
	return w.XYZToENU(ref, xyz)
}

func (w wgs84) XYZToENU(ref LLA, e XYZ) ENU {
	var (
		lambda = ref.Latitude.Radians().F64()
		phi    = ref.Longitude.Radians().F64()

		sinLambda = math.Sin(lambda)
		cosLambda = math.Cos(lambda)
		sinPhi    = math.Sin(phi)
		cosPhi    = math.Cos(phi)

		xref = w.LLAToXYZ(ref)

		xd = (e.X - xref.X).F64()
		yd = (e.Y - xref.Y).F64()
		zd = (e.Z - xref.Z).F64()
	)

	return ENU{
		East:  Meters(-sinPhi*xd + cosPhi*yd),
		North: Meters(-cosPhi*sinLambda*xd - sinLambda*sinPhi*yd + cosLambda*zd),
		Up:    Meters(cosLambda*cosPhi*xd + cosLambda*sinPhi*yd + sinLambda*zd),
	}
}

func (w wgs84) ENUToXYZ(ref LLA, e ENU) XYZ {
	var (
		lambda = ref.Latitude.Radians().F64()
		phi    = ref.Longitude.Radians().F64()
		s      = math.Sin(lambda)
		n      = wgs84A / math.Sqrt(1-wgs84ESq*s*s)

		sinLambda = math.Sin(lambda)
		cosLambda = math.Cos(lambda)

		sinPhi = math.Sin(phi)
		cosPhi = math.Cos(phi)

		x0 = (ref.Altitude.F64() + n) * cosLambda * cosPhi
		y0 = (ref.Altitude.F64() + n) * cosLambda * sinPhi
		z0 = (ref.Altitude.F64() + (1-wgs84ESq)*n) * sinLambda

		east  = e.East.F64()
		north = e.North.F64()
		up    = e.Up.F64()

		xd = -sinPhi*east - cosPhi*sinLambda*north + cosLambda*cosPhi*up
		yd = cosPhi*east - sinLambda*sinPhi*north + cosLambda*sinPhi*up
		zd = cosLambda*north + sinLambda*up

		x = xd + x0
		y = yd + y0
		z = zd + z0
	)

	return XYZ{
		X: Meters(x),
		Y: Meters(y),
		Z: Meters(z),
	}
}

// vim: foldmethod=marker
