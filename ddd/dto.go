package ddd

type PagedItems[T any] struct {
	Items   []T  `json:"items"`
	HasNext bool `json:"hasNext"`
}

func NewPagedItems[T any](items []T, hasNext bool) PagedItems[T] {
	return PagedItems[T]{
		HasNext: hasNext,
		Items:   items,
	}
}

type CheckFieldReason string

const (
	// ReasonOk is used when the value is ok
	ReasonOk CheckFieldReason = "ok"
	// ReasonExistedValue is used when the value is existed
	ReasonExistedValue CheckFieldReason = "value_existed"
	// ReasonRequired is used when the field is required
	ReasonRequired CheckFieldReason = "required"
	// ReasonInvalidField is used when the field is not valid
	ReasonInvalidField CheckFieldReason = "invalid_field"
	// ReasonNotSupport is used when the field is not supported
	ReasonNotSupport CheckFieldReason = "not_support"
)

type CheckFieldResult struct {
	Result  CheckFieldReason `json:"reason"`
	Message string           `json:"message"`
}

func CheckFieldOk() *CheckFieldResult {
	return &CheckFieldResult{Result: ReasonOk}
}
