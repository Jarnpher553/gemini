package dto

// Response响应类
type Response struct {
	ErrCode   int         `json:"code"`
	ErrMsg    string      `json:"msg"`
	Timestamp int64       `json:"timestamp"`
	Success   bool        `json:"success"`
	Data      interface{} `json:"data"`
}

// TokenOut 带token的data
type TokenOut struct {
	Token string      `json:"token"`
	Datum interface{} `json:"datum"`
}

type PagedOut struct {
	TotalCount int         `json:"total_count"`
	PerCount   int         `json:"per_count"`
	PageNum    int         `json:"page_num"`
	QueryCount int         `json:"query_count"`
	Rows       interface{} `json:"rows"`
}

type ListOut struct {
	Rows interface{} `json:"rows"`
}
