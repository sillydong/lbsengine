package core

import (
	"fmt"
	"github.com/sillydong/lbsengine/types"
	"math/rand"
	"testing"
)

var indexer *Indexer

func init() {
	indexer = &Indexer{}
	indexer.Init(&types.IndexerOptions{
		RedisHost:   "127.0.0.1:6379",
		RedisDb:     3,
		HashSize:    1000,
		GeoShard:    5,
		GeoPrecious: 5,
	})
}

func TestAddDocument(t *testing.T) {
	doc := &types.IndexedDocument{
		DocId:     1,
		Latitude:  40.7137674,
		Longitude: -73.9525142,
		Fields: map[string]string{
			"a": "b",
			"c": "d",
		},
	}
	indexer.Add(doc)
}

func TestRemoveDocument(t *testing.T) {
	indexer.Remove(1)
}

func TestStore(t *testing.T) {
	strs, err := indexer.client.HVals("h_dr5rt_1").Result()
	if err != nil {
		t.Error(err)
	}
	docs := make([]types.IndexedDocument, 0)
	for _, str := range strs {
		document := types.IndexedDocument{}
		document.UnmarshalMsg([]byte(str))
		docs = append(docs, document)
	}
	fmt.Printf("%+v\n", docs)
}

func BenchmarkStore(b *testing.B) {
	str, _ := indexer.client.HGet("h_dr5rt_1", "1").Result()
	for i := 0; i < b.N; i++ {
		document := types.IndexedDocument{}
		document.UnmarshalMsg(tobytes(str))
	}
}

func TestSearch(t *testing.T) {
	option := &types.SearchOptions{
		Refresh:   false,
		OrderDesc: false,
		Accuracy:  types.STANDARD,
		Circles:   1,
		Excepts: map[uint64]bool{
			1: true,
		},
		Filter: func(doc types.IndexedDocument) bool {
			fmt.Println("filtering")
			if doc.Fields.(map[string]interface{})["a"] == "b" {
				return true
			}
			return false
		},
	}
	docs, num := indexer.Search(false, "h_dr5rt_1", 40.7137674, -73.9525142, option)
	fmt.Println(num)
	if num > 0 {
		for _, doc := range docs {
			fmt.Printf("%+v\n", doc)
		}
	}
}

func BenchmarkSearch(b *testing.B) {
	option := &types.SearchOptions{
		Refresh:   false,
		OrderDesc: false,
		Accuracy:  types.MEITUAN,
		Circles:   1,
		//Excepts: map[uint64]bool{
		//	1: true,
		//},
		//Filter: func(doc types.IndexedDocument) bool {
		//	fmt.Println("filtering")
		//	if doc.Fields.(map[string]interface{})["a"] == "b" {
		//		return true
		//	}
		//	return false
		//},
	}
	for i := 0; i < b.N; i++ {
		indexer.Search(false, "h_dr5rt_1", 40.7137674, -73.9525142, option)
	}
}

func TestIdShard(t *testing.T) {
	for i := 0; i < 10; i++ {
		x := rand.Int()
		fmt.Printf("%v -> %v\n", x, indexer.hashshard(uint64(x)))
	}
}

func TestGeoShard(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Printf("%v -> %v\n", i, indexer.geoshard("asdfg", uint64(i)))
	}
}
