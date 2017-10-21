package engine

import (
	"github.com/sillydong/lbsengine/core"
	"github.com/sillydong/lbsengine/types"
)

type Engine struct {
	option *types.EngineOptions

	cachers  []*core.Cacher
	indexers []*core.Indexer

	indexerAddChannels    []chan *indexerAddRequest
	indexerRemoveChannels []chan *indexerRemoveRequest
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

	e.indexerAddChannels = make([]chan *indexerAddRequest,e.option.NumShards)
	e.indexerRemoveChannels = make([]chan *indexerRemoveRequest,e.option.NumShards)
	e.indexerSearchChannels = make([]chan *indexerSearchRequest,e.option.NumShards)

	for i:=0;i<e.option.NumShards;i++{
		//初始化缓存器
		e.cachers[i]=&core.Cacher{}
		//初始化索引器
		e.indexers[i]=&core.Indexer{}
		e.indexers[i].Init(e.option.IndexerOption)
		//初始化channel
		e.indexerAddChannels[i]=make(chan *indexerAddRequest,e.option.AddBuffer)
		e.indexerRemoveChannels[i]=make(chan *indexerRemoveRequest,e.option.RemoveBuffer)
		e.indexerSearchChannels[i]=make(chan *indexerSearchRequest,e.option.SearchBuffer)

		go indexerAddWorker(i)
		go indexerRemoveWorker(i)
		for k:=0;k<e.option.SearchWorkerThreads;k++{
			go indexerSearchWorker(i)
		}
	}
}

func (e *Engine) Add(doc *types.IndexedDocument) {

}

func (e *Engine) Remove() {

}

func (e *Engine) Search() {

}
