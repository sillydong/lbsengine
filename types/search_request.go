package types

import "time"

//搜索请求
type SearchRequest struct {
	Latitude     float64
	Longitude    float64
	CountOnly    bool
	Offset       int
	Limit        int
	SearchOption *SearchOptions //可留空，使用引擎默认参数
}

//可不设置的搜索参数
type SearchOptions struct {
	Refresh   bool
	OrderDesc bool
	Timeout   time.Duration
	Accuracy  int             //计算进度
	Circles   int             //圈数，默认1，不扩散
	Excepts   map[uint64]bool //排除指定ID
	Filter    func(doc IndexedDocument) bool
}

func (o *SearchOptions) Init() {
	if o.Accuracy == 0 {
		o.Accuracy = STANDARD
	}
	if o.Circles == 0 {
		o.Circles = 1
	}
}

const (
	_        = iota
	STANDARD //传统计算方法
	MEITUAN  //美团开放计算方法
	IMPROVED //优化的计算方法
)
