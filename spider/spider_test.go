package spider

import (
	"fmt"
	"testing"

	"github.com/3zheng/railgun/PoolAndAgent"
)

func TestGetGaodeData(t *testing.T) {
	//初始化并连接mysql数据库
	pDBProcess := PoolAndAgent.CreateADODatabase("root:123456@tcp(localhost:3306)/gotest?charset=utf8")
	pDBProcess.InitDB()

	for i := 1; i <= 40; i++ {
		str := fmt.Sprintf("%d", i)
		arrData := GetPOIData(str)
		fmt.Println("长度=", len(arrData), arrData)
		t.Log(arrData)

		for _, value := range arrData {
			sqlExpress := fmt.Sprintf("insert into poi_data(id, name, location, longitude, latitude) values( '%s', '%s', '%s', %f, %f)", value.ID, value.Name, value.Location, value.Coordinate.Longitude, value.Coordinate.Latitude)
			pDBProcess.WriteToDB(sqlExpress)
		}
	}

}
