package distanceMeasure

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

const TEST_DATA_NUM = 5000000

var arrPoint1 [TEST_DATA_NUM]*EarthCoordinate                                          //测试数据组1
var arrPoint2 [TEST_DATA_NUM]*EarthCoordinate                                          //测试数据组2
var precision float64 = 0.1                                                            //随机数生成的精度大小
var benchmarkPoint EarthCoordinate = EarthCoordinate{Longitude: 121.0, Latitude: 31.0} //基准值
var resultStardard, resultQuick, resultQuick2 [TEST_DATA_NUM]float64                   //各个算法的测距后得到的结果

func init() {
	//使用随机数生成测试数据
	var tmpPoint EarthCoordinate
	now := time.Now().Unix()
	rand.Seed(now)

	for i := 0; i < TEST_DATA_NUM; i++ {
		//随机生成1000000组经纬度坐标
		//以1000000为除数求余，然后把得到的余数再与除以1000000作为benchmarkPoint的小数部分
		tmpPoint.Latitude = benchmarkPoint.Latitude + float64(rand.Intn(1000000))/1000000*precision
		tmpPoint.Longitude = benchmarkPoint.Longitude + float64(rand.Intn(1000000))/1000000*precision
		arrPoint1[i] = &tmpPoint
		tmpPoint.Latitude = benchmarkPoint.Latitude + float64(rand.Intn(1000000))/1000000*precision
		tmpPoint.Longitude = benchmarkPoint.Longitude + float64(rand.Intn(1000000))/1000000*precision
		arrPoint2[i] = &tmpPoint
	}
}

//标准球体算法
func TestStardard(t *testing.T) {
	measure := GetInstance()
	//计算标准三维距离算法的时间
	before := time.Now().UnixNano()
	for i := 0; i < TEST_DATA_NUM; i++ {
		resultStardard[i] = measure.MeasureByStardardMethod(arrPoint1[i], arrPoint2[i])
	}
	after := time.Now().UnixNano()
	t.Logf("Haversine公式标准球体算法耗时：%E ns\n", float64(after-before))
}

//没有预设本地经纬度坐标的快速测距算法
func TestQuickWithoutSetLocation(t *testing.T) {
	measure := GetInstance()
	//通过三维近似为二维计算距离算法的时间
	before := time.Now().UnixNano()
	for i := 0; i < TEST_DATA_NUM; i++ {
		resultQuick[i] = measure.MeasureByQuickMethodWithoutLocation(arrPoint1[i], arrPoint2[i])
	}
	after := time.Now().UnixNano()
	t.Logf("平面近似二维距离算法（勾股定理，未预设城市经纬度）的耗时：%E ns\n", float64(after-before))
}

func TestQuick(t *testing.T) {
	measure := GetInstance()
	//先来个错误的示范
	if _, err := measure.MeasureByQuickMethod(arrPoint1[0], arrPoint2[0]); err != nil {
		t.Log(err)
	}

	//通过三维近似为二维计算距离算法的时间
	before := time.Now().UnixNano()
	measure.SetLocalEarthCoordinate(benchmarkPoint, "上海")
	for i := 0; i < TEST_DATA_NUM; i++ {
		if value, err := measure.MeasureByQuickMethod(arrPoint1[i], arrPoint2[i]); err != nil {
			t.Error("出错了")
		} else {
			resultQuick2[i] = value
		}
	}
	after := time.Now().UnixNano()
	t.Logf("平面近似二维距离算法（勾股定理,已预设城市经纬度）的耗时：%E ns\n", float64(after-before))
}

//计算标准差
func TestCompare(t *testing.T) {
	var variance, variance2 float64 = 0.0, 0.0
	for i := 0; i < TEST_DATA_NUM; i++ {
		//	log.Printf("第%d组数据: EarthPointPair = %v, stardard = %f, quick = %f, quick2 = %f\n",
		//		i, arrPoint[i], resultStardard[i], resultQuick[i], resultQuick[i])
		variance += math.Pow(resultQuick[i]-resultStardard[i], 2)
		variance2 += math.Pow(resultQuick2[i]-resultStardard[i], 2)
	}
	t.Logf("精度为%f度时，换算以米为单位时表示搜索范围距离在%f米内时, 标准差组1为%f米,组2为%f米\n",
		precision, DIST_PER_DEGREE*precision, math.Sqrt(variance/TEST_DATA_NUM), math.Sqrt(variance2/TEST_DATA_NUM))
	return
}
