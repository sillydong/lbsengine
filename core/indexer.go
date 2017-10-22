package core

import (
	"github.com/go-redis/redis"
	"github.com/mmcloughlin/geohash"
	"github.com/sillydong/lbsengine/distanceMeasure"
	"github.com/sillydong/lbsengine/types"
	"log"
	"strconv"
	"unsafe"
)

type Indexer struct {
	option  *types.IndexerOptions
	client  *redis.Client
	measure *distanceMeasure.DistanceMeasure
}

func (i *Indexer) Init(option *types.IndexerOptions) {
	i.option = option
	i.client = redis.NewClient(&redis.Options{
		Addr:     option.RedisHost,
		DB:       option.RedisDb,
		Password: option.RedisPassword,
	})
	err := i.client.Ping().Err()
	if err != nil {
		log.Fatal(err)
	}
	//距离计算
	i.measure = distanceMeasure.GetInstance()
	if option.CenterLongitude != 0.0 && option.CenterLatitude != 0.0 {
		i.measure.SetLocalEarthCoordinate(&distanceMeasure.EarthCoordinate{Latitude: option.CenterLatitude, Longitude: option.CenterLongitude}, option.Location)
	}
}

func (i *Indexer) Add(doc *types.IndexedDocument) {
	sdocid := strconv.FormatUint(doc.DocId, 10)
	hashs := i.hashshard(doc.DocId)
	pip := i.client.Pipeline()
	//删除旧数据
	oldhash, _ := i.client.HGet(hashs, sdocid).Result()
	if len(oldhash) > 0 {
		pip.HDel(i.geoshard(oldhash, doc.DocId), sdocid)
	}
	newhash := geohash.EncodeWithPrecision(doc.Latitude, doc.Longitude, i.option.GeoPrecious)
	pip.HSet(hashs, sdocid, newhash)
	pip.HSet(i.geoshard(newhash, doc.DocId), sdocid, doc)
	_, err := pip.Exec()
	if err != nil {
		log.Print(err)
	}
}

func (i *Indexer) Remove(docid uint64) {
	sdocid := strconv.FormatUint(docid, 10)
	hashs := i.hashshard(docid)
	hash, _ := i.client.HGet(hashs, sdocid).Result()
	if len(hash) > 0 {
		geos := i.geoshard(hash, docid)
		pip := i.client.Pipeline()
		pip.HDel(hashs, sdocid)
		pip.HDel(geos, sdocid)
		_, err := pip.Exec()
		if err != nil {
			log.Print(err)
		}
	}
}

func (i *Indexer) Search(countonly bool, hash string, latitude, longitude float64, options *types.SearchOptions) (docs types.ScoredDocuments, count int) {
	strs := i.client.HVals(hash).Val()
	if len(strs) > 0 {
		if countonly {
			//仅计算数量
			if options.Excepts == nil && options.Filter == nil {
				count += len(strs)
			} else {
				for _, str := range strs {
					document := types.IndexedDocument{}
					document.UnmarshalMsg(tobytes(str))
					if document.DocId != 0 {
						//判断是否排除
						if options.Excepts == nil {
							//判断是否过滤
							if options.Filter == nil || options.Filter(document) {
								count++
							}
						} else {
							if _, ok := options.Excepts[document.DocId]; !ok {
								if options.Filter == nil || options.Filter(document) {
									count++
								}
							}
						}
					}
				}
			}
		} else {
			//需要数据
			for _, str := range strs {
				document := types.IndexedDocument{}
				document.UnmarshalMsg(tobytes(str))
				if document.DocId != 0 {
					//判断是否排除
					if options.Excepts == nil {
						//判断是否过滤
						if options.Filter == nil || options.Filter(document) {
							doc := types.ScoredDocument{
								DocId:    document.DocId,
								Distance: i.distance(options.Accuracy, document.Latitude, document.Longitude, latitude, longitude),
							}
							docs = append(docs, &doc)
							count++
						}
					} else {
						if _, ok := options.Excepts[document.DocId]; !ok {
							if options.Filter == nil || options.Filter(document) {
								doc := types.ScoredDocument{
									DocId:    document.DocId,
									Distance: i.distance(options.Accuracy, document.Latitude, document.Longitude, latitude, longitude),
								}
								docs = append(docs, &doc)
								count++
							}
						}
					}
				}
			}
		}
	}
	return
}

//距离计算
func (i *Indexer) distance(accuracy int, alatitude, alongitude, blatitude, blongitude float64) float64 {
	//fmt.Printf("%+v - %+v ------ %+v - %+v",alatitude,alongitude,blatitude,blongitude)
	a := &distanceMeasure.EarthCoordinate{Latitude: alatitude, Longitude: alongitude}
	b := &distanceMeasure.EarthCoordinate{Latitude: blatitude, Longitude: blongitude}
	switch accuracy {
	case types.STANDARD:
		return i.measure.MeasureByStardardMethod(a, b)
	case types.MEITUAN:
		return i.measure.MeasureByQuickMethodWithoutLocation(a, b)
	case types.IMPROVED:
		distance, err := i.measure.MeasureByQuickMethod(a, b)
		if err != nil {
			return i.measure.MeasureByQuickMethodWithoutLocation(a, b)
		}
		return distance
	}
	return 0.0
}

func (i *Indexer) hashshard(docid uint64) string {
	return "g_" + strconv.FormatUint((uint64(docid/i.option.HashSize)+1)*i.option.HashSize, 10)
}

func (i *Indexer) geoshard(hash string, docid uint64) string {
	return "h_" + hash + "_" + strconv.FormatUint(docid-docid/i.option.GeoShard*i.option.GeoShard, 10)
}

func tobytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
