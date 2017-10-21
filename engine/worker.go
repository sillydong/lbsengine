package engine

import (
	"github.com/sillydong/lbsengine/types"
)

type indexerSearchRequest struct {
	countonly            bool
	hash                 string
	latitude             float64
	longitude            float64
	option               *types.SearchOptions
	indexerReturnChannel chan *indexerSearchResponse
}

type indexerSearchResponse struct {
	docs  types.ScoredDocuments
	count int
}

func (e *Engine) indexerAddWorker(shard int) {
	request := <-e.indexerAddChannels[shard]
	e.indexer.Add(request)
}

func (e *Engine) indexerRemoveWorker(shard int) {
	request := <-e.indexerRemoveChannels[shard]
	e.indexer.Remove(request)
}

func (e *Engine) indexerSearchWorker(shard int) {
	request := <-e.indexerSearchChannels[shard]
	docs, count := e.indexer.Search(request.countonly, request.hash, request.latitude, request.longitude, request.option)
	request.indexerReturnChannel <- &indexerSearchResponse{docs: docs, count: count}
}
