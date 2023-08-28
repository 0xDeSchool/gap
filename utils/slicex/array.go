package slicex

func Unique[TSource any](source []TSource, equals func(t1 TSource, t2 TSource) bool) []TSource {
	newArr := make([]TSource, 0)
	for i := 0; i < len(source); i++ {
		repeat := false
		for j := i + 1; j < len(source); j++ {
			if equals(source[i], source[j]) {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, source[i])
		}
	}
	return newArr
}

func Map[TSource any, TResult any](source []TSource, selector func(*TSource) TResult) []TResult {
	if source == nil {
		return nil
	}
	if selector == nil {
		panic("parameter selector is nil")
	}
	result := make([]TResult, len(source))
	for i := 0; i < len(source); i++ {
		v := &source[i]
		result[i] = selector(v)
	}
	return result
}

func Sum[TSource any](source []TSource, selector func(*TSource) int) int {
	if source == nil {
		panic("parameter source is nil")
	}
	if selector == nil {
		panic("parameter selector is nil")
	}
	count := 0
	for i := 0; i < len(source); i++ {
		count += selector(&source[i])
	}
	return count
}

func IndexOf[TSource any](source []TSource, predicate func(*TSource) bool) int {
	if source == nil {
		panic("parameter source is nil")
	}
	if predicate == nil {
		panic("parameter selector is nil")
	}
	for i := 0; i < len(source); i++ {
		v := &source[i]
		if predicate(v) {
			return i
		}
	}
	return -1
}

func Contains[TSource comparable](source []TSource, item *TSource) bool {
	if source == nil {
		return false
	}

	for i := 0; i < len(source); i++ {
		if source[i] == *item {
			return true
		}
	}
	return false
}

func MapMany[TSource any, TResult any](source []TSource, selector func(*TSource) []TResult) []TResult {
	if source == nil {
		panic("parameter source is nil")
	}
	if selector == nil {
		panic("parameter selector is nil")
	}
	result := make([]TResult, 0)
	for i := 0; i < len(source); i++ {
		v := &source[i]
		result = append(result, selector(v)...)
	}
	return result
}

func First[TSource any](source []TSource, predicate func(*TSource) bool) *TSource {
	if source == nil {
		panic("parameter source is nil")
	}
	if predicate == nil {
		panic("parameter predicate is nil")
	}
	for i := 0; i < len(source); i++ {
		v := &source[i]
		if predicate(v) {
			return v
		}
	}
	return nil
}

func Filter[TSource any](source []TSource, predicate func(*TSource) bool) []TSource {
	if source == nil {
		return nil
	}
	if predicate == nil {
		panic("parameter predicate is nil")
	}
	result := make([]TSource, 0)
	for i := 0; i < len(source); i++ {
		v := &source[i]
		if predicate(v) {
			result = append(result, *v)
		}
	}
	return result
}

func FilterCount[TSource any](source []TSource, predicate func(*TSource) bool) int {
	if source == nil {
		return 0
	}
	if predicate == nil {
		panic("parameter predicate is nil")
	}
	count := 0
	for i := 0; i < len(source); i++ {
		v := &source[i]
		if predicate(v) {
			count += 1
		}
	}
	return count
}

func FilterMap[TSource any, TResult any](source []TSource, selector func(*TSource) (TResult, bool)) []TResult {
	if source == nil {
		panic("parameter source is nil")
	}
	if selector == nil {
		panic("parameter selector is nil")
	}
	result := make([]TResult, 0)
	for i := 0; i < len(source); i++ {
		v := &source[i]
		if r, ok := selector(v); ok {
			result = append(result, r)
		}
	}
	return result
}

func ToMap[TSource any, TKey comparable](source []TSource, keySelector func(*TSource) TKey) map[TKey]*TSource {
	if source == nil {
		panic("parameter source is nil")
	}
	if keySelector == nil {
		panic("parameter keySelector is nil")
	}
	result := make(map[TKey]*TSource)
	for i := 0; i < len(source); i++ {
		v := &source[i]
		result[keySelector(v)] = v
	}
	return result
}
func ToMapPtr[TSource any, TKey comparable](source []*TSource, keySelector func(*TSource) TKey) map[TKey]*TSource {
	if source == nil {
		panic("parameter source is nil")
	}
	if keySelector == nil {
		panic("parameter keySelector is nil")
	}
	result := make(map[TKey]*TSource)
	for i := 0; i < len(source); i++ {
		v := source[i]
		result[keySelector(v)] = v
	}
	return result
}

func Values[TSource any, TKey comparable](source map[TKey]TSource) []TSource {
	if source == nil {
		panic("parameter source is nil")
	}
	result := make([]TSource, 0, len(source))
	for _, v := range source {
		result = append(result, v)
	}
	return result
}

func Keys[TSource any, TKey comparable](source map[TKey]TSource) []TKey {
	if source == nil {
		return nil
	}
	result := make([]TKey, 0, len(source))
	for k := range source {
		result = append(result, k)
	}
	return result
}

func Distinct[TSource any, TKey comparable](source []TSource, keySelector func(*TSource) TKey) []TKey {
	if source == nil {
		panic("parameter source is nil")
	}
	if keySelector == nil {
		panic("parameter keySelector is nil")
	}
	result := make([]TKey, 0)
	dict := make(map[TKey]struct{})
	for i := 0; i < len(source); i++ {
		k := keySelector(&source[i])
		if _, ok := dict[k]; !ok {
			result = append(result, k)
			dict[k] = struct{}{}
		}
	}
	return result
}

func DistinctBy[TSource any, TKey comparable](source []TSource, keySelector func(*TSource) TKey) []TKey {
	if source == nil {
		panic("parameter source is nil")
	}
	if keySelector == nil {
		panic("parameter keySelector is nil")
	}
	result := make([]TKey, 0)
	dict := make(map[TKey]struct{})
	for i := 0; i < len(source); i++ {
		k := keySelector(&source[i])
		if _, ok := dict[k]; !ok {
			result = append(result, k)
			dict[k] = struct{}{}
		}
	}
	return result
}

func DistinctMany[TSource any, TKey comparable](source []TSource, keySelector func(*TSource) []TKey) []TKey {
	if source == nil {
		return nil
	}
	if keySelector == nil {
		panic("parameter keySelector is nil")
	}
	result := make([]TKey, 0)
	dict := make(map[TKey]struct{})
	for i := 0; i < len(source); i++ {
		keys := keySelector(&source[i])
		for _, k := range keys {
			if _, ok := dict[k]; !ok {
				result = append(result, k)
				dict[k] = struct{}{}
			}
		}
	}
	return result
}

func ToHashSet[TSource any, TKey comparable](source []TSource, keySelector func(*TSource) TKey) map[TKey]struct{} {
	result := make(map[TKey]struct{})
	for i := range source {
		result[keySelector(&source[i])] = struct{}{}
	}
	return result
}

func GroupBy[TSource any, TKey comparable](source []TSource, keySelector func(*TSource) TKey) map[TKey][]*TSource {
	if source == nil {
		panic("parameter source is nil")
	}
	if keySelector == nil {
		panic("parameter keySelector is nil")
	}
	result := make(map[TKey][]*TSource)
	for i := range source {
		k := keySelector(&source[i])
		result[k] = append(result[k], &source[i])
	}
	return result
}

// Split split array by max length of each segment
func Split[TSource any](source []TSource, maxSegmentLength int, maxChars int) [][]TSource {
	if source == nil {
		panic("parameter source is nil")
	}
	if maxSegmentLength <= 0 {
		panic("parameter maxSegmentLength must be greater than 0")
	}
	result := make([][]TSource, 0)
	for i := 0; i < len(source); i += maxSegmentLength {
		end := i + maxSegmentLength
		if end > len(source) {
			end = len(source)
		}
		result = append(result, source[i:end])
	}
	return result
}
