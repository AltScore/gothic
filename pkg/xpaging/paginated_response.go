package xpaging

type PaginatedResponse[C any] struct {
	Items         []C `json:"items"`
	PagingOptions `json:",inline"`
	Total         int64 `json:"total"`
}
