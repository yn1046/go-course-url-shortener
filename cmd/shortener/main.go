package main

import (
	"crypto/sha1"
	"io"
	"net/http"
	"strings"
)

type MemStorage struct {
	storage map[string]string
}

var globalStorage MemStorage = MemStorage{storage: make(map[string]string)}

func mainHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		reqData, _ := io.ReadAll(req.Body)

		urlString := string(reqData)
		h := sha1.New()
		io.WriteString(h, urlString)
		hashString := string(h.Sum((nil)))[:8]

		globalStorage.storage[hashString] = urlString

		body := "http://localhost:8080/" + urlString
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(body))
	} else if req.Method == http.MethodGet {
		hash := strings.TrimPrefix(req.URL.Path, "/")
		res.Write([]byte(globalStorage.storage[hash]))
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, mainHandler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
