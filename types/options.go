package types

import (
	"runtime"
	"time"
)

type EngineOptions struct {
	NumShards           int //channel分片
	AddBuffer           int //add channel长度
	RemoveBuffer        int //remove channel长度
	SearchBuffer        int //channel长度
	SearchWorkerThreads int //每个搜索channel worker数量

	DefaultSearchOption *SearchOptions  //默认搜索配置
	IndexerOption       *IndexerOptions //索引器配置项
}

func (o *EngineOptions) Init() {
	if o.NumShards == 0 {
		o.NumShards = runtime.NumCPU()
	}
	if o.AddBuffer == 0 {
		o.AddBuffer = runtime.NumCPU()
	}
	if o.RemoveBuffer == 0 {
		o.RemoveBuffer = runtime.NumCPU()
	}
	if o.SearchBuffer == 0 {
		o.SearchBuffer = 10 * runtime.NumCPU()
	}
	if o.SearchWorkerThreads == 0 {
		o.SearchWorkerThreads = 9 //9个neighbor格子
	}
	if o.DefaultSearchOption == nil {
		o.DefaultSearchOption = &SearchOptions{
			Refresh:   false,
			OrderDesc: false,
			Timeout:   2 * time.Second,
			Accuracy:  STANDARD,
			Circles:   1,
			Excepts:   nil,
			Filter:    nil,
		}
	}
	if o.IndexerOption == nil {
		o.IndexerOption = &IndexerOptions{
			RedisHost:     "127.0.0.1:6379",
			RedisPassword: "",
			RedisDb:       3,
			HashSize:      1000,
			GeoShard:      5,
			GeoPrecious:   5,
		}
	}
}

type IndexerOptions struct {
	RedisHost       string
	RedisPassword   string
	RedisDb         int
	HashSize        uint64  //hash分片大小
	GeoShard        uint64  //GEOHASH分片大小
	GeoPrecious     uint    //GEOHASH位数
	CenterLatitude  float64 //城市中心纬度
	CenterLongitude float64 //城市中心经度
	Location        string  //城市
}
