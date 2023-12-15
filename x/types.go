package x

type PageParam struct {
	Page         int64 `json:"page"`         // 页码, 从1开始
	PageSize     int64 `json:"pageSize"`     // 每页大小
	IncludeTotal bool  `json:"includeTotal"` // 是否包含总数
	IncludedId   any   `json:"id"`           // 包含指定Id的数据，用于分页查询，这时候不需要传入page
}

func (p *PageParam) Skip() int64 {
	return (p.Page - 1) * p.PageSize
}

func (p *PageParam) Limit() int64 {
	return p.PageSize
}

// NewPageParam new page param
func NewPageParam(page, pageSize int64) *PageParam {
	return &PageParam{
		Page:     page,
		PageSize: pageSize,
	}
}

func PageParamWithId(includedId any, limit int64) *PageParam {
	return &PageParam{
		IncludedId: includedId,
		PageSize:   limit,
	}
}

type PageAndSort struct {
	PageParam `json:",inline"` // 页码, 从1开始
	Sort      string           `json:"sort"` // 排序字段, 例如: -createdAt, +updatedAt
}

func (p *PageAndSort) Copy() *PageAndSort {
	return &PageAndSort{
		PageParam: PageParam{
			Page:     p.Page,
			PageSize: p.PageSize,
		},
		Sort: p.Sort,
	}
}

func PageLimit(limit int) *PageAndSort {
	return &PageAndSort{
		PageParam: PageParam{
			Page:     1,
			PageSize: int64(limit),
		},
	}
}

// NewPageAndSort new page and sort
func NewPageAndSort(page, pageSize int64, sort string) *PageAndSort {
	return &PageAndSort{
		PageParam: PageParam{
			Page:     page,
			PageSize: pageSize,
		},
		Sort: sort,
	}
}

func NewPageAndSortWithId(includedId any, pageSize int64, sort string) *PageAndSort {
	return &PageAndSort{
		PageParam: PageParam{
			IncludedId: includedId,
			PageSize:   pageSize,
		},
		Sort: sort,
	}
}

// IsDesc is sorted by desc
func (p *PageAndSort) IsDesc() bool {
	return len(p.Sort) > 0 && p.Sort[0] == '-' || p.Sort[0] == '!'
}

func IsSortDesc(sort string) bool {
	return len(sort) > 0 && sort[0] == '-' || sort[0] == '!'
}

type PagedResult[T any] struct {
	Total   int64 `json:"total"`
	Data    []T   `json:"data"`
	HasMore bool  `json:"hasMore"`
	Page    int64 `json:"page"`
}
