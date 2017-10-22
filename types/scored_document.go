package types

//排序的struct
type ScoredDocument struct {
	Distance float64
	DocId    uint64
	Model    interface{} //为外围组装数据提供
}

//排序
type ScoredDocuments []*ScoredDocument

func (docs ScoredDocuments) Len() int {
	return len(docs)
}

func (docs ScoredDocuments) Swap(i, j int) {
	docs[i], docs[j] = docs[j], docs[i]
}

func (docs ScoredDocuments) Less(i, j int) bool {
	return docs[i].Distance < docs[j].Distance
}
