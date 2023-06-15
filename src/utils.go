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
