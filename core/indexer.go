package core

import (
	"github.com/go-redis/redis"
	"github.com/sillydong/lbsengine/types"
	"log"
	"strconv"
	"github.com/mmcloughlin/geohash"
)

type Indexer struct {
	option *types.IndexerOptions
	client *redis.Client
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
}

func (i *Indexer) Add(doc *types.IndexedDocument) {
	hash := geohash.EncodeWithPrecision(doc.Latitude, doc.Longitude, i.option.GeoPrecious)
	sdocid := strconv.FormatUint(doc.DocId, 10)
	pip := i.client.Pipeline()
	pip.HSet(i.hashshard(doc.DocId), sdocid, hash)
	pip.HSet(i.geoshard(hash, doc.DocId), sdocid, doc)
	_, err := pip.Exec()
	if err != nil {
		log.Print(err)
	}
}

func (i *Indexer) Remove(docid uint64) {
	sdocid := strconv.FormatUint(docid, 10)
	hashs := i.hashshard(docid)
	hash := i.client.HGet(hashs, sdocid).String()
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

func (i *Indexer) Search(latitude,longitude float64,options *types.SearchOptions) (docs types.ScoredDocuments,count int){
	return nil,0
}

func (i *Indexer) hashshard(docid uint64) string {
	return "g_" + strconv.FormatUint((uint64(docid/i.option.HashSize)+1)*i.option.HashSize, 10)
}

func (i *Indexer) geoshard(hash string, docid uint64) string {
	return "h_" + hash + "_" + strconv.FormatUint(docid-docid/i.option.GeoShard*i.option.GeoShard, 10)
}
