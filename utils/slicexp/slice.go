package slicexp

func First[TSource any](source []TSource, predicate func(TSource) bool) (TSource, bool) {
	if source == nil {
		panic("parameter source is nil")
	}
	if predicate == nil {
		panic("parameter predicate is nil")
	}
	for i := 0; i < len(source); i++ {
		if predicate(source[i]) {
			return source[i], true
		}
	}
	var res TSource
	return res, false
}

func Map[TSource any, TResult any](source []*TSource, selector func(*TSource) TResult) []TResult {
	if source == nil {
		return nil
	}
	if selector == nil {
		panic("parameter selector is nil")
	}
	result := make([]TResult, len(source))
	for i := 0; i < len(source); i++ {
		result[i] = selector(source[i])
	}
	return result
}

func MapMany[TSource any, TResult any](source []*TSource, selector func(*TSource) []TResult) []TResult {
	if source == nil {
		panic("parameter source is nil")
	}
	if selector == nil {
		panic("parameter selector is nil")
	}
	result := make([]TResult, 0)
	for i := 0; i < len(source); i++ {
		result = append(result, selector(source[i])...)
	}
	return result
}

func Filter[TSource any](source []*TSource, predicate func(*TSource) bool) []*TSource {
	if source == nil {
		return nil
	}
	if predicate == nil {
		panic("parameter predicate is nil")
	}
	result := make([]*TSource, 0)
	for _, v := range source {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

func ToMap[TSource any, TKey comparable](source []*TSource, keySelector func(*TSource) TKey) map[TKey]*TSource {
	if source == nil {
		panic("parameter source is nil")
	}
	if keySelector == nil {
		panic("parameter keySelector is nil")
	}
	result := make(map[TKey]*TSource)
	for _, v := range source {
		result[keySelector(v)] = v
	}
	return result
}

func Distinct[TSource comparable](source []*TSource, keySelector func(*TSource) TSource) []TSource {
	if source == nil {
		return nil
	}
	if keySelector == nil {
		panic("parameter keySelector is nil")
	}
	result := make([]TSource, 0)
	dict := make(map[TSource]struct{})
	for i := 0; i < len(source); i++ {
		k := keySelector(source[i])
		if _, ok := dict[k]; !ok {
			result = append(result, k)
			dict[k] = struct{}{}
		}
	}
	return result
}

func DistinctBy[TSource any, TKey comparable](source []*TSource, keySelector func(*TSource) TKey) []TKey {
	if source == nil {
		return nil
	}
	if keySelector == nil {
		panic("parameter keySelector is nil")
	}
	result := make([]TKey, 0)
	dict := make(map[TKey]struct{})
	for i := 0; i < len(source); i++ {
		k := keySelector(source[i])
		if _, ok := dict[k]; !ok {
			result = append(result, k)
			dict[k] = struct{}{}
		}
	}
	return result
}

func ToHashSet[TSource any, TKey comparable](source []*TSource, keySelector func(*TSource) TKey) map[TKey]struct{} {
	result := make(map[TKey]struct{})
	for i := range source {
		result[keySelector(source[i])] = struct{}{}
	}
	return result
}

func CountBy[TSource any](source []*TSource, predicate func(*TSource) bool) int {
	if source == nil {
		return 0
	}
	if predicate == nil {
		panic("parameter predicate is nil")
	}
	count := 0
	for i := 0; i < len(source); i++ {
		if predicate(source[i]) {
			count += 1
		}
	}
	return count
}

func FirstOrNil[TSource any](source []*TSource, predicate func(*TSource) bool) *TSource {
	if source == nil {
		return nil
	}
	if predicate == nil {
		panic("parameter predicate is nil")
	}
	for _, v := range source {
		if predicate(v) {
			return v
		}
	}
	return nil
}

func Contains[TSource comparable](source []*TSource, item *TSource) bool {
	if source == nil {
		return false
	}

	for i := 0; i < len(source); i++ {
		if source[i] == item {
			return true
		}
	}
	return false
}

func Sum[TSource any](source []*TSource, selector func(*TSource) int) int {
	if source == nil {
		return 0
	}
	if selector == nil {
		panic("parameter selector is nil")
	}
	count := 0
	for i := 0; i < len(source); i++ {
		count += selector(source[i])
	}
	return count
}

func GroupBy[TSource any, TKey comparable](source []TSource, keySelector func(TSource) TKey) map[TKey][]TSource {
	if source == nil {
		panic("parameter source is nil")
	}
	if keySelector == nil {
		panic("parameter keySelector is nil")
	}
	result := make(map[TKey][]TSource)
	for i := range source {
		k := keySelector(source[i])
		result[k] = append(result[k], source[i])
	}
	return result
}
