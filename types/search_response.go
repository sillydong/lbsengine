package types

//搜索结果返回
type SearchResponse struct{
	Docs ScoredDocuments //排序好的结果
	Timeout bool //是否超时
	Count int64 //数量
}
