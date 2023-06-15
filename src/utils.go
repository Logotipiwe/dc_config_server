package main

import (
	"fmt"
	"net/http"
)

func getScheme(r *http.Request) string {
	if r.TLS == nil {
		return "http"
	}
	return "https"
}

func getCurrHost(r *http.Request) string {
	return fmt.Sprintf("%s://%s", getScheme(r), r.Host)
}
