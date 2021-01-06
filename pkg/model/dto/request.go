package dto

// PagedIn 分页输入
type PagedIn struct {
	PageNum  int `json:"page_num" form:"page_num" binding:"gte=1"`
	PerCount int `json:"per_count" form:"per_count" binding:"gte=1"`
}
