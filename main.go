package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func hello(w http.ResponseWriter, req *http.Request) {
	host, _ := os.Hostname()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "not set"
	}
	io.WriteString(w, fmt.Sprintf("[v3]Hello, Kubernetes, From host: %s, DB_URL: %s\n", host, dbURL))
}

func main() {
	http.HandleFunc("/", hello)
	http.ListenAndServe(":3000", nil)
}
