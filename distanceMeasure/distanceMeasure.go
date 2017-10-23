package distanceMeasure

/*
使用本package的DistanceMeasure结构体建议使用GetInstance和CreateNewMeasure来获取对象，
而不建议自己直接new。因为如果自己new可能会在切换城市的时候忘记设置本地的location
GetInstance和CreateNewMeasure的区别是：
1.如果一个进程同时只为一个城市或一个区（比如上海这么人口密集的城市可以一个区一个基准坐标，这样误差也会小）服务，
那么建议使用GetInstance通常只需要在启动时以读入配置文件的配置项的方式来设置就行了。
2.如果一个进程同时为多个城市或者多个区服务，那么建议使用CreateNewMeasure，这样调用SetLocalEarthCoordinate
设置本地城市坐标时才不会互相覆盖和影响
MeasureByStardardMethod的数学推导：http://blog.csdn.net/liminlu0314/article/details/8553926
MeasureByQuickMethodWithoutLocation的参考链接：https://tech.meituan.com/lucene-distance.html
*/

import (
	"fmt"
	"math"
)

const RADIUS = 6378137.0                      //地球半径，单位米
const DIST_PER_DEGREE = math.Pi * 35433.88889 // πr/180° 解释为每度所表示的实际长度

var localMeasure *DistanceMeasure = nil

type DistanceMeasure struct {
	IsSetLocation bool             //是否设置了本地基准经纬度
	Benchmark     *EarthCoordinate //基准坐标
	cosLatitude   float64          //math.Cos(Benchmark.latitude) 基准维度的cos值
	cityName      string           //城市名，可以不填，打印使用的
	IsFirstUse    bool             //是否是第一次使用，在GetInstance使用后置为ture
}

type MeasureError struct {
}

func (e MeasureError) Error() string {
	return "尚未输入本地城市的基准经纬度坐标"
}

type EarthCoordinate struct {
	Longitude float64 //经度坐标 东经为正数，西经为负数
	Latitude  float64 //纬度坐标 北纬为正数，南纬为负数
}

//返回单例
func GetInstance() *DistanceMeasure {
	if localMeasure == nil {
		localMeasure = new(DistanceMeasure)
		localMeasure.IsSetLocation = false
		localMeasure.Benchmark = &EarthCoordinate{0.0, 0.0}
	}
	localMeasure.IsFirstUse = true
	return localMeasure
}

//创建一个新的DistanceMeasure对象
func CreateNewMeasure() *DistanceMeasure{
	measure := new(DistanceMeasure)
	measure.IsSetLocation = false
	measure.Benchmark = &EarthCoordinate{0.0, 0.0}
	measure.IsFirstUse = true
	return measure
}

//设置本地的基准坐标，可以取一个城市的市中心，EarthCoordinate必须要输入所在城市的经纬度坐标，name值（城市名）可以为空
func (this *DistanceMeasure) SetLocalEarthCoordinate(location *EarthCoordinate, name string) {
	this.IsSetLocation = true                                                //已设置坐标
	this.Benchmark = location                                                //用于QuickMethod的本地经纬度坐标
	this.cityName = name                                                     //城市名，打印使用
	this.cosLatitude = math.Cos(this.ChangeAngleToRadian(location.Latitude)) //基准维度的cos值
	fmt.Println("设置了城市坐标，城市名：", name, "经纬度坐标：", location)
}

func (this *DistanceMeasure) ChangeAngleToRadian(angle float64) float64 {
	return angle / 180.0 * math.Pi
}

//标准球体测距算法，基准算法，基于球面模型来处理的（立体几何），即Haversine公式
func (this *DistanceMeasure) MeasureByStardardMethod(pt1, pt2 *EarthCoordinate) float64 {
	//先把角度转成弧度
	lon1 := this.ChangeAngleToRadian(pt1.Longitude)
	lat1 := this.ChangeAngleToRadian(pt1.Latitude)
	lon2 := this.ChangeAngleToRadian(pt2.Longitude)
	lat2 := this.ChangeAngleToRadian(pt2.Latitude)
	//因为cos余弦函数是偶函数，所以lon2-lon1还是lon1-lon2都一样
	//得到距离的弧度值
	radDist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(lon2-lon1))
	return radDist * RADIUS //乘以地球半径，返回实际长度
}

//未提前设置本地城市经纬度坐标的快速测距算法
func (this *DistanceMeasure) MeasureByQuickMethodWithoutLocation(pt1, pt2 *EarthCoordinate) float64 {
	diffLat := (pt1.Latitude - pt2.Latitude) //纬度差的实际距离，单位m
	radLat := this.ChangeAngleToRadian((pt1.Latitude + pt2.Latitude) / 2)
	diffLog := math.Cos(radLat) * (pt1.Longitude - pt2.Longitude) //经度差的实际距离，单位m,需要乘以cos(所在纬度),因为纬度越高时经度差的表示的实际距离越短
	//根据勾股定理
	dist := DIST_PER_DEGREE * math.Sqrt(diffLat*diffLat+diffLog*diffLog)
	return dist
}

//已提前设置本地城市经纬度坐标的快速测距算法
func (this *DistanceMeasure) MeasureByQuickMethod(pt1, pt2 *EarthCoordinate) (float64, error) {
	//先判断是否已经输入了本地基准经纬度坐标

	if this.IsSetLocation == false {
		return 0.0, MeasureError{}
	}
	if this.IsFirstUse {
		this.IsFirstUse = false
		//第一次调用这个函数的时候打印城市名进行提示
		var strLon, strLat string
		if this.Benchmark.Longitude >= 0 {
			strLon = fmt.Sprintf("东经%f°", this.Benchmark.Longitude)
		} else {
			strLon = fmt.Sprintf("西经%f°", -this.Benchmark.Longitude)
		}
		if this.Benchmark.Latitude >= 0 {
			strLat = fmt.Sprintf("北纬%f°", this.Benchmark.Latitude)
		} else {
			strLat = fmt.Sprintf("南纬%f°", -this.Benchmark.Latitude)
		}
		fmt.Printf("本次快速测距算法设置的城市为%s,输入的基准经纬度坐标为[%s, %s]\n", this.cityName, strLon, strLat)
		fmt.Println("如果当前保存的城市与你期望的城市不符，请重新调用SetLocalEarthCoordinate来设置")
	}

	diffLat := (pt1.Latitude - pt2.Latitude)                      //纬度差的实际距离，单位m
	diffLog := this.cosLatitude * (pt1.Longitude - pt2.Longitude) //经度差的实际距离，单位m,需要乘以cos(所在纬度),因为纬度越高时经度差的表示的实际距离越短
	distance := DIST_PER_DEGREE * math.Sqrt(diffLat*diffLat+diffLog*diffLog)
	return distance, nil
}
