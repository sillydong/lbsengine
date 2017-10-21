package dbtest

import (
	"math"
)

const RADIUS 1

type EarthCoordinate struct {
	longitude float64 //经度坐标 东经为正数，西经为负数
	latitude  float64 //纬度坐标 北纬为正数，南纬为负数
}

func ChangeAngleToRadian(angle float64) float64{
	return angle / 180.0 * math.Pi
}

func TestEarthPoint(pt1, pt2 EarthCoordinate){
	x1 := RADIUS * math.Cos(ChangeAngleToRadian(pt1.latitude)) * math.Cos(ChangeAngleToRadian(pt1.longitude))
	y1 := RADIUS * math.Cos(ChangeAngleToRadian(pt1.latitude)) * math.Sin(ChangeAngleToRadian(pt1.longitude))
	z1 := RADIUS * math.Sin(ChangeAngleToRadian(pt1.latitude))
	
	x2 := RADIUS * math.Cos(ChangeAngleToRadian(pt2.latitude)) * math.Cos(ChangeAngleToRadian(pt2.longitude))
	y2 := RADIUS * math.Cos(ChangeAngleToRadian(pt2.latitude)) * math.Sin(ChangeAngleToRadian(pt2.longitude))
	z2 := RADIUS * math.Sin(ChangeAngleToRadian(pt2.latitude))
}
