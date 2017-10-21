package engine

import (
	"github.com/sillydong/lbsengine/core"
	"github.com/sillydong/lbsengine/types"
	"fmt"
	"unsafe"
	"github.com/huichen/murmur"
)

type Engine struct {
	option *types.EngineOptions

	cachers  []*core.Cacher
	indexers []*core.Indexer

	indexerAddChannels    []chan *types.IndexedDocument
	indexerRemoveChannels []chan uint64
	indexerSearchChannels []chan *indexerSearchRequest
}

func(e *Engine)Init(option *types.EngineOptions){
	if option == nil{
		option = &types.EngineOptions{}
	}
	option.Init()
	e.option=option

	e.cachers = make([]*core.Cacher,e.option.NumShards)
	e.indexers = make([]*core.Indexer,e.option.NumShards)

	e.indexerAddChannels = make([]chan *types.IndexedDocument,e.option.NumShards)
	e.indexerRemoveChannels = make([]chan uint64,e.option.NumShards)
	e.indexerSearchChannels = make([]chan *indexerSearchRequest,e.option.NumShards)

	for i:=0;i<e.option.NumShards;i++{
		//初始化缓存器
		e.cachers[i]=&core.Cacher{}
		//初始化索引器
		e.indexers[i]=&core.Indexer{}
		e.indexers[i].Init(e.option.IndexerOption)
		//初始化channel
		e.indexerAddChannels[i]=make(chan *types.IndexedDocument,e.option.AddBuffer)
		e.indexerRemoveChannels[i]=make(chan uint64,e.option.RemoveBuffer)
		e.indexerSearchChannels[i]=make(chan *indexerSearchRequest,e.option.SearchBuffer)

		go e.indexerAddWorker(i)
		go e.indexerRemoveWorker(i)
		for k:=0;k<e.option.SearchWorkerThreads;k++{
			go e.indexerSearchWorker(i)
		}
	}
}

func (e *Engine) Add(doc *types.IndexedDocument) {
	shard := e.shard(doc.DocId)
	e.indexerAddChannels[shard]<-doc
}

func (e *Engine) Remove(docid uint64) {
	shard := e.shard(docid)
	e.indexerRemoveChannels[shard]<-docid
}

func (e *Engine) Search() {

}

func(e *Engine)shard(docid uint64)int{
	s := fmt.Sprintf("%d",docid)
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	hash := murmur.Murmur3(*(*[]byte)(unsafe.Pointer(&h)))
	return int(hash - hash/e.option.NumShards*e.option.NumShards)
}
