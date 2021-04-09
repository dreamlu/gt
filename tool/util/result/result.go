package result

// Pager info
type Pager struct {
	ClientPage int64 `json:"client_page"` // page num
	EveryPage  int64 `json:"every_page"`  // page size
	TotalNum   int64 `json:"total_num"`   // data num
}
