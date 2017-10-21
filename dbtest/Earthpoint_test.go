package dbtest

import (
	"fmt"
	"math"
)

const RADIUS = 1.0

type EarthCoordinate struct {
	longitude float64 //经度坐标 东经为正数，西经为负数
	latitude  float64 //纬度坐标 北纬为正数，南纬为负数
}

func ChangeAngleToRadian(angle float64) float64 {
	return angle / 180.0 * math.Pi
}

func CalEarthPoint(pt1, pt2 EarthCoordinate) {
	x1 := RADIUS * math.Cos(ChangeAngleToRadian(pt1.latitude)) * math.Cos(ChangeAngleToRadian(pt1.longitude))
	y1 := RADIUS * math.Cos(ChangeAngleToRadian(pt1.latitude)) * math.Sin(ChangeAngleToRadian(pt1.longitude))
	z1 := RADIUS * math.Sin(ChangeAngleToRadian(pt1.latitude))

	x2 := RADIUS * math.Cos(ChangeAngleToRadian(pt2.latitude)) * math.Cos(ChangeAngleToRadian(pt2.longitude))
	y2 := RADIUS * math.Cos(ChangeAngleToRadian(pt2.latitude)) * math.Sin(ChangeAngleToRadian(pt2.longitude))
	z2 := RADIUS * math.Sin(ChangeAngleToRadian(pt2.latitude))

	distance := math.Sqrt(math.Pow(x1-x2, 2.0) + math.Pow(y1-y2, 2.0) + math.Pow(z1-z2, 2.0))
	fmt.Printf("x1 = %f, y1 = %f, z1 = %f\nx2 = %f, y2 = %f, z2 = %f\ndistance = %f", x1, y1, z1, x2, y2, z2, distance)
}

func Test_EarthPoint() {
	p1a, p2a := EarthCoordinate{longitude: 100.0, latitude: 20.0}, EarthCoordinate{longitude: 130.0, latitude: 40.0}
	p1b, p2b := EarthCoordinate{longitude: 170.0, latitude: 20.0}, EarthCoordinate{longitude: -160.0, latitude: 40.0}
	p1c, p2c := EarthCoordinate{longitude: 170.0, latitude: -20.0}, EarthCoordinate{longitude: -160.0, latitude: -40.0}
	CalEarthPoint(p1a, p2a)
	CalEarthPoint(p1b, p2b)
	CalEarthPoint(p1c, p2c)
}
