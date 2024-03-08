package main

func toMap[K string, V any](arr []V, groupingFunc func(val V) K) map[K]V {
	result := make(map[K]V)
	for _, v := range arr {
		key := groupingFunc(v)
		result[key] = v
	}
	return result
}

func Filter[T any](slice []T, filterFunc func(T) bool) []T {
	result := make([]T, 0)
	for _, elem := range slice {
		if filterFunc(elem) {
			result = append(result, elem)
		}
	}
	return result
}
