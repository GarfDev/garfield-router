package main

import (
	_ "embed"
	"net/http"
)

//go:embed ui-usage.html
var usageHTML []byte

func handleUsage(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	_, _ = w.Write(usageHTML)
}
