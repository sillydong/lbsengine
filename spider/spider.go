package spider

import (
	"crypto/md5"
	"fmt"
	"github.com/sillydong/lbsengine/distanceMeasure"
	"io/ioutil"
	"net/http"
	"sort"
)

type PoiData struct {
	ID         string                          //唯一标识符
	Name       string                          //兴趣点的名字
	Location   string                          //字符串形式的坐标
	Coordinate distanceMeasure.EarthCoordinate //EarthCoordinate形式的坐标
}

type URL struct {
	UrlHead  string            //url的开始部分，类似http://restapi.amap.com/v3/place/text?
	MapParam map[string]string //url的参数部分
}

func (this *URL) Init(urlHead string) {
	this.MapParam = make(map[string]string, 50)
	this.UrlHead = urlHead
}

func (this *URL) AddParam(key, value string) {
	this.MapParam[key] = value
}

func (this *URL) GetFinalURL() string {
	str := this.UrlHead
	for key, value := range this.MapParam {
		str += "&" + key + "=" + value
	}
	return str
}

//根据已有的参数获取数字签名的MD5值,privateKey为私钥
func (this *URL) GetMD5Sign(privateKey string) string {
	keys := make([]string, 0)
	for key, _ := range this.MapParam {
		keys = append(keys, key)
	}
	//升序排序
	sort.Strings(keys)
	//根据排序后的结果取key-value值
	var sign string
	length := len(keys)
	for i, key := range keys {
		elem, _ := this.MapParam[key]
		sign += key + "=" + elem
		if i < length-1 {
			//不是最后一个的话要+&
			sign += "&"
		} else {
			//最后一个则直接+私钥
			sign += privateKey
		}
	}
	//进行MD5加密
	fmt.Println("被加密的MD5串：", sign)
	byteMD5 := md5.Sum([]byte(sign))
	sigMD5 := fmt.Sprintf("%x", byteMD5)
	return sigMD5
}

func GetPOIData(pageId string) []PoiData {
	//var url string = "http://restapi.amap.com/v3/place/text?&keywords=超市&city=shanghai&key=7ffc2abae565cdb9302ebaef2e45c572&extensions=all&sig="
	pUrl := new(URL)
	pUrl.Init("http://restapi.amap.com/v3/place/text?")
	pUrl.AddParam("keywords", "超市")
	pUrl.AddParam("city", "shanghai")
	pUrl.AddParam("key", "7ffc2abae565cdb9302ebaef2e45c572")
	pUrl.AddParam("extensions", "all")
	pUrl.AddParam("page", pageId)
	//获取数字签名的MD5值
	sig := pUrl.GetMD5Sign("f66d435f57cb361630f2110a97aa4fd9")
	pUrl.AddParam("sig", sig)
	finalUrl := pUrl.GetFinalURL()
	fmt.Println("final url: ", finalUrl)
	resp, err := http.Get(finalUrl)
	if err != nil {
		fmt.Println("post错误:", err)
		return nil
	}
	fmt.Println("resp结果值：", resp)

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("返回的body为：", string(body[:]))
	return ReadFromJson(body)
}
