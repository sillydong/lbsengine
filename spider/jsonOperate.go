package spider

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func ReadFromJson(body []byte) (datas []PoiData) {
	datas = make([]PoiData, 0)
	var dat map[string]interface{}
	dat = make(map[string]interface{})
	err := json.Unmarshal(body, &dat)
	if err != nil {
		fmt.Println("json解析错误:", err)
		return
	}

	if suggestion, ok := dat["suggestion"]; ok {
		fmt.Println("type of suggestion: ", reflect.TypeOf(suggestion))
	}

	if pois, ok := dat["pois"]; ok {
		fmt.Println("type of pois: ", reflect.TypeOf(pois))
		switch data := pois.(type) {
		case []interface{}:
			fmt.Println("POI数组长度:", len(data))
			for _, elem := range data {
				switch elem := elem.(type) {
				case map[string]interface{}:
					name, _ := elem["name"]
					location, _ := elem["location"]
					fmt.Println(name, location)
					data := PoiData{}
					switch name := name.(type) {
					case string:
						data.Name = name
					}
					switch location := location.(type) {
					case string:
						data.Location = location
						coord := strings.Split(location, ",")
						data.Coordinate.Longitude, _ = strconv.ParseFloat(coord[0], 64)
						data.Coordinate.Latitude, _ = strconv.ParseFloat(coord[1], 64)
					}
					datas = append(datas, data)
				}
			}
		}
	}

	fmt.Println("datas =", datas)
	return
}
