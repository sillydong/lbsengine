package engine

import (
	"fmt"
	"github.com/huichen/murmur"
	"github.com/sillydong/lbsengine/core"
	"github.com/sillydong/lbsengine/types"
	"sort"
	"strconv"
	"time"
	"unsafe"
)

type Engine struct {
	option *types.EngineOptions

	cacher  *core.Cacher
	indexer *core.Indexer

	indexerAddChannels    []chan *types.IndexedDocument
	indexerRemoveChannels []chan uint64
	indexerSearchChannels []chan *indexerSearchRequest
}

func (e *Engine) Init(option *types.EngineOptions) {
	if option == nil {
		option = &types.EngineOptions{}
	}
	option.Init()
	e.option = option

	//初始化缓存器
	e.cacher = &core.Cacher{}
	e.cacher.Init()
	//初始化索引器
	e.indexer = &core.Indexer{}
	e.indexer.Init(e.option.IndexerOption)

	e.indexerAddChannels = make([]chan *types.IndexedDocument, e.option.NumShards)
	e.indexerRemoveChannels = make([]chan uint64, e.option.NumShards)
	e.indexerSearchChannels = make([]chan *indexerSearchRequest, e.option.NumShards)

	for i := 0; i < int(e.option.NumShards); i++ {
		//初始化channel
		e.indexerAddChannels[i] = make(chan *types.IndexedDocument, e.option.AddBuffer)
		e.indexerRemoveChannels[i] = make(chan uint64, e.option.RemoveBuffer)
		e.indexerSearchChannels[i] = make(chan *indexerSearchRequest, e.option.SearchBuffer)

		go e.indexerAddWorker(i)
		go e.indexerRemoveWorker(i)
		for k := 0; k < e.option.SearchWorkerThreads; k++ {
			go e.indexerSearchWorker(i)
		}
	}
}

func (e *Engine) Add(doc *types.IndexedDocument) {
	shard := e.shardid(doc.DocId)
	e.indexerAddChannels[shard] <- doc
}

func (e *Engine) Remove(docid uint64) {
	shard := e.shardid(docid)
	e.indexerRemoveChannels[shard] <- docid
}

func (e *Engine) Search(request *types.SearchRequest) *types.SearchResponse {
	content := fmt.Sprintf("%v", request)
	hash := murmur.Murmur3(tobytes(content))

	shard := int(hash - hash/uint32(e.option.NumShards)*uint32(e.option.NumShards))

	cachekey := fmt.Sprintf("%v", hash)

	if request.SearchOption == nil {
		request.SearchOption = e.option.DefaultSearchOption
	}
	request.SearchOption.Init()

	//是否刷新缓存
	if !request.SearchOption.Refresh {
		//从缓存取数据
		docs, count := e.cacher.Get(cachekey, request.Offset, request.Limit)
		if docs != nil && count > 0 {
			return &types.SearchResponse{Docs: docs, Count: count, Timeout: false}
		}
	}

	//获取neighbour
	neighbours := core.LoopNeighbours(request.Latitude, request.Longitude, e.option.IndexerOption.GeoPrecious, request.SearchOption.Circles)
	loops := len(neighbours) * int(e.option.IndexerOption.GeoShard) //geohash分片
	indexerReturnChannel := make(chan *indexerSearchResponse, loops)

	//下发任务
	for _, geo := range neighbours {
		for i := 0; i < int(e.option.IndexerOption.GeoShard); i++ {
			geoshard := "h_" + geo + "_" + strconv.Itoa(i)
			e.indexerSearchChannels[shard] <- &indexerSearchRequest{
				countonly:            request.CountOnly,
				hash:                 geoshard,
				latitude:             request.Latitude,
				longitude:            request.Longitude,
				option:               request.SearchOption,
				indexerReturnChannel: indexerReturnChannel,
			}
		}
	}

	//整理数据
	docs := types.ScoredDocuments{}
	count := 0
	istimeout := false
	if request.SearchOption.Timeout > 0 {
		deadline := time.Now().Add(request.SearchOption.Timeout)
		for i := 0; i < loops; i++ {
			select {
			case lresponse := <-indexerReturnChannel:
				if lresponse.count == 0 {
					continue
				}
				if !request.CountOnly {
					docs = append(docs, lresponse.docs...)
				}
				count += lresponse.count
			case <-time.After(deadline.Sub(time.Now())):
				istimeout = true
				break
			}

		}
	} else {
		for i := 0; i < loops; i++ {
			lresponse := <-indexerReturnChannel
			if lresponse.count == 0 {
				continue
			}
			if !request.CountOnly {
				docs = append(docs, lresponse.docs...)
			}
			count += lresponse.count
		}
	}

	//最终排序
	if !request.CountOnly {
		if request.SearchOption.OrderDesc {
			sort.Sort(sort.Reverse(docs))
		} else {
			sort.Sort(docs)
		}
	}

	//写缓存
	if len(docs) > 0 {
		e.cacher.Set(cachekey, docs)
	}

	//拼返回数据
	response := &types.SearchResponse{}
	response.Count = count
	if !request.CountOnly {
		start := request.Offset
		stop := request.Offset + request.Limit
		if start > count {
			response.Docs = nil
		} else if stop > count {
			response.Docs = docs[start:]
		} else {
			response.Docs = docs[start:stop]
		}
	} else {
		response.Docs = nil
	}
	response.Timeout = istimeout

	return response

}

func (e *Engine) shardid(docid uint64) int {
	content := fmt.Sprintf("%d", docid)
	hash := murmur.Murmur3(tobytes(content))
	return int(hash - hash/uint32(e.option.NumShards)*uint32(e.option.NumShards))
}

func tobytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
