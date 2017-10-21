package distanceMeasure

import (
	"fmt"
	"math"
)

type EarthCoordinate struct {
	longitude float64 //经度坐标 东经为正数，西经为负数
	latitude  float64 //纬度坐标 北纬为正数，南纬为负数
}

func MeasureByStardardMethod(pt1, pt2 EarthCoordinate) float64 {
	if math.Abs(pt1.longitude-pt2.longitude) > 180 {

	}
}
