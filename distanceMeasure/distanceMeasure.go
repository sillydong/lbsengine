package distanceMeasure

/*
使用本package的DistanceMeasure结构体建议是用GetInstance来获取对象，而不建议自己直接new。
因为如果自己new可能会在切换城市的时候忘记设置本地的location
*/

import (
	"fmt"
	"math"
)

const RADIUS = 6378137.0                      //地球半径，单位米
const DIST_PER_DEGREE = math.Pi * 35433.88889 // πr/180° 解释为每度所表示的实际长度

var localMeasure *DistanceMeasure = nil

type DistanceMeasure struct {
	IsSetLocation bool            //是否设置了本地基准经纬度
	Benchmark     EarthCoordinate //基准坐标
	cosLatitude   float64         //math.Cos(Benchmark.latitude) 基准维度的cos值
	cityName      string          //城市名，可以不填，打印使用的
	IsFirstUse    bool            //是否是第一次使用，在GetInstance使用后置为ture
}

type MeasureError struct {
}

func (e MeasureError) Error() string {
	return "尚未输入本地城市的基准经纬度坐标"
}

type EarthCoordinate struct {
	longitude float64 //经度坐标 东经为正数，西经为负数
	latitude  float64 //纬度坐标 北纬为正数，南纬为负数
}

//返回单例
func GetInstance() *DistanceMeasure {
	if localMeasure == nil {
		localMeasure = new(DistanceMeasure)
		localMeasure.IsSetLocation = false
		localMeasure.Benchmark = EarthCoordinate{0.0, 0.0}
	}
	localMeasure.IsFirstUse = true
	return localMeasure
}

//设置本地的基准坐标，可以取一个城市的市中心，EarthCoordinate必须要输入所在城市的经纬度坐标，name值（城市名）可以为空
func (this *DistanceMeasure) SetLocalEarthCoordinate(location EarthCoordinate, name string) {
	this.IsSetLocation = true                                                //已设置坐标
	this.Benchmark = location                                                //用于QuickMethod的本地经纬度坐标
	this.cityName = name                                                     //城市名，打印使用
	this.cosLatitude = math.Cos(this.ChangeAngleToRadian(location.latitude)) //基准维度的cos值
}

func (this *DistanceMeasure) ChangeAngleToRadian(angle float64) float64 {
	return angle / 180.0 * math.Pi
}

func (this *DistanceMeasure) MeasureByStardardMethod(pt1, pt2 EarthCoordinate) float64 {
	//先把角度转成弧度
	lon1 := this.ChangeAngleToRadian(pt1.longitude)
	lat1 := this.ChangeAngleToRadian(pt1.latitude)
	lon2 := this.ChangeAngleToRadian(pt2.longitude)
	lat2 := this.ChangeAngleToRadian(pt2.latitude)
	//因为cos余弦函数是偶函数，所以lon2-lon1还是lon1-lon2都一样
	//得到距离的弧度值
	radDist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(lon2-lon1))
	return radDist * RADIUS //乘以地球半径，返回实际长度
}

//未提前设置本地城市经纬度坐标的快速测距算法
func (this *DistanceMeasure) MeasureByQuickMethodNotSetLocation(pt1, pt2 EarthCoordinate) float64 {
	diffLat := (pt1.latitude - pt2.latitude) //纬度差的实际距离，单位m
	radLat := this.ChangeAngleToRadian((pt1.latitude + pt2.latitude) / 2)
	diffLog := math.Cos(radLat) * (pt1.longitude - pt2.longitude) //经度差的实际距离，单位m,需要乘以cos(所在纬度),因为纬度越高时经度差的表示的实际距离越短
	//根据勾股定理
	dist := DIST_PER_DEGREE * math.Sqrt(diffLat*diffLat+diffLog*diffLog)
	return dist
}

//已提前设置本地城市经纬度坐标的快速测距算法
func (this *DistanceMeasure) MeasureByQuickMethod(pt1, pt2 EarthCoordinate) (float64, error) {
	//先判断是否已经输入了本地基准经纬度坐标
	if this.IsSetLocation == false {
		return 0.0, MeasureError{}
	}
	if this.IsFirstUse {
		//第一次调用这个函数的时候打印城市名进行提示
		var strLon, strLat string
		if this.Benchmark.longitude >= 0 {
			strLon = fmt.Sprintf("东经%f°", this.Benchmark.longitude)
		} else {
			strLon = fmt.Sprintf("西经%f°", -this.Benchmark.longitude)
		}
		if this.Benchmark.latitude >= 0 {
			strLat = fmt.Sprintf("北纬%f°", this.Benchmark.longitude)
		} else {
			strLat = fmt.Sprintf("南纬%f°", -this.Benchmark.longitude)
		}
		fmt.Printf("本次快速测距算法设置的城市为%s,输入的基准经纬度坐标为[%s, %s]", this.cityName, strLon, strLat)
		fmt.Println("如果当前保存的城市与你期望的城市不符，请重新调用SetLocalEarthCoordinate来设置")
	}

	diffLat := (pt1.latitude - pt2.latitude)                      //纬度差的实际距离，单位m
	diffLog := this.cosLatitude * (pt1.longitude - pt2.longitude) //经度差的实际距离，单位m,需要乘以cos(所在纬度),因为纬度越高时经度差的表示的实际距离越短
	distance := DIST_PER_DEGREE * math.Sqrt(diffLat*diffLat+diffLog*diffLog)
	return distance, nil
}
