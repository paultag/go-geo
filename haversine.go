/* {{{ Copyright (c) Paul R. Tagliamonte <paultag@gmail.com>, 2015
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE. }}} */

package geo

import (
	"fmt"
	"math"
)

var earthRadiusMeters float64 = 6371000

// HaversineDistance will return the haversine (great-circle) distance between
// two Lat/Lon points as if one were to traverse the surface of the earth.
//
// Because distance between two points uses the radius of the Earth, this means
// HaversineDistance only really makes sense at an elevation of 0. Trying to
// determine the great-circle distance that also accounts for an elevation
// difference is conceptually annoying (and hard to debug) -- and I generally
// don't know when you'd *actually* want that metric, since if you care about
// the curve of the earth, the elevation difference ought to not be the biggest
// component of that distance.
//
// As a result of my stubborn insistence to drive down errors due to API misuse,
// this function will return an error if either of the provided geo.LLA structs
// have an Altitude other than 0.
func HaversineDistance(origin, position LLA) (Meters, error) {
	if origin.Altitude != 0 || position.Altitude != 0 {
		return 0, fmt.Errorf("geo.HaversineDistance: Altitude must be 0")
	}

	var (
		originLon = origin.Longitude.Radians().F64()
		originLat = origin.Latitude.Radians().F64()

		positionLon = position.Longitude.Radians().F64()
		positionLat = position.Latitude.Radians().F64()

		deltaLon = originLon - positionLon
		deltaLat = originLat - positionLat
	)

	a := math.Pow(math.Sin(deltaLat/2), 2) + math.Cos(originLat)*math.Cos(positionLat)*math.Pow(math.Sin(deltaLon/2), 2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return Meters(earthRadiusMeters * c), nil
}

// vim: foldmethod=marker
