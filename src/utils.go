package main

import (
	"fmt"
	"os"
)

func getScheme() string {
	return os.Getenv("OUTER_SCHEME")
}

func getCurrHost() string {
	return fmt.Sprintf("%s://%s:%s",
		getScheme(), os.Getenv("OUTER_HOST"), os.Getenv("OUTER_PORT"))
}

func getSubpath() string {
	return os.Getenv("SUBPATH")
}

func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}

func toMap[K string, V any](arr []V, groupingFunc func(val V) K) map[K]V {
	result := make(map[K]V)
	for _, v := range arr {
		key := groupingFunc(v)
		result[key] = v
	}
	return result
}
