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

// Package geo contains helpers to deal with the conversion and manipulation
// of locations.
package geo

import (
	"math"
)

// Meters represents the SI unit of measure, Meters.
type Meters float64

// F64 will return the value as a float64. Doing "value.F64()" is the same
// as doing "float64(value)", except this can be a bit more clean at times.
func (m Meters) F64() float64 {
	return float64(m)
}

// Degrees represents an angular measurement, in Degrees. This type is for
// two main reasons -- firstly, to enforce (on a type level) that the user is
// aware that the angle measurements must be in Degrees, and second, to bind
// a conversion helper to Radians.
type Degrees float64

// F64 will return the value as a float64. Doing "value.F64()" is the same
// as doing "float64(value)", except this can be a bit more clean at times.
func (d Degrees) F64() float64 {
	return float64(d)
}

// Radians returns the Angle, but in terms of Radians.
func (d Degrees) Radians() Radians {
	return Radians(math.Pi / 180 * d)
}

// Radians represents an angular measurement, in Radians. This type is for
// two main reasons -- firstly, to enforce (on a type level) that the user is
// aware that the angle measurements must be in Radians, and second, to bind
// a conversion helper to Degrees.
type Radians float64

// F64 will return the value as a float64. Doing "value.F64()" is the same
// as doing "float64(value)", except this can be a bit more clean at times.
func (r Radians) F64() float64 {
	return float64(r)
}

// Degrees returns the Angle, but in terms of Degrees.
func (r Radians) Degrees() Degrees {
	return Degrees(180 / math.Pi * r)
}

// CoordinateSystem is a system to map locations (usually Latitude and
// Longitude, in the form of LLA objects) to absolute points in space.
//
// Different systems have different measurements of Earth's surface, and
// a Latitude / Longitude must be understood within its CoordinateSystem,
// or significant errors can be introduced.
type CoordinateSystem interface {

	// LLAToXYZ will take a LLA inside this coordinate system and return that
	// in absolute XYZ space.
	LLAToXYZ(LLA) XYZ

	// XYZToLLA will take an absolute XYZ and return an LLA.
	XYZToLLA(XYZ) LLA

	// XYZToENU will take an absolute XYZ and return that on the ENU tangent
	// plane at the provided LLA.
	XYZToENU(LLA, XYZ) ENU

	// ENUToXYZ will take a relative ENU and translate that into an absolute
	// XYZ given the tangent plane at the reference LLA.
	ENUToXYZ(LLA, ENU) XYZ

	// LLAToENU will return the ENU relative to the first LLA of the second LLA,
	// returned in the ENU plane.
	LLAToENU(LLA, LLA) ENU
}

// AER represents an Azimuth, Elevation, Range measurement.
//
// Azimuth/Elevation (or Az/El) is a common way of locating objects measured
// at a specific location (for instance, from a RADAR)
//
// This is a *relative* and *angular* measure.
type AER struct {
	Azimuth   Degrees
	Elevation Degrees
	Range     Meters
}

// LLA or Latitude, Longitude, Altitude, is a location somewhere around Earth.
//
// This is a *absolute* and *angular* measure.
type LLA struct {
	Latitude  Degrees
	Longitude Degrees
	Altitude  Meters
}

// XYZ is the earth-centric XYZ point system LLA locations can be turned into
// points on the Earth's ellipsoid, but plotted using cartesian coordinates
// relative to Earth, rather than angular LLA measurements.
//
// This is a *absolute* and *cartesian* measure.
type XYZ struct {
	X Meters
	Y Meters
	Z Meters
}

// ENU is East, North, Up in Meters. These measures are in the local
// tangent plane, which is to say, increasing "North" will get further and
// further away from the Earth's surface (well, unless there's a mountain
// range in front of you).
//
// This is a *relative* and *cartesian* measure.
type ENU struct {
	East  Meters
	North Meters
	Up    Meters
}

// ENU will translate the AER angular vector into 3D space as an ENU.
func (aed AER) ENU() ENU {
	var (
		sinTheta = math.Sin(aed.Elevation.Radians().F64())
		cosTheta = math.Cos(aed.Elevation.Radians().F64())
		sinPhi   = math.Sin(aed.Azimuth.Radians().F64())
		cosPhi   = math.Cos(aed.Azimuth.Radians().F64())
	)

	return ENU{
		East:  aed.Range * Meters(sinPhi*cosTheta),
		North: aed.Range * Meters(cosPhi*sinTheta),
		Up:    aed.Range * Meters(sinTheta),
	}
}

// AER will convert the 3D ENU point, and return it as an angular AER vector.
func (enu ENU) AER() AER {
	var r = Meters(math.Sqrt((enu.East*enu.East + enu.North*enu.North).F64()))

	return AER{
		Azimuth:   Degrees(math.Atan2(enu.East.F64(), enu.North.F64())),
		Elevation: Degrees(math.Atan2(enu.Up.F64(), r.F64())),
		Range:     Meters(math.Sqrt((r*r + enu.Up*enu.Up).F64())),
	}
}

// vim: foldmethod=marker
