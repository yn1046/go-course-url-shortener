package main

import (
	"fmt"
	"io"
	"net/http"
)

func mainPage(res http.ResponseWriter, req *http.Request) {
	body := fmt.Sprintf("Method: %s\r\n", req.Method)

	reqData, _ := io.ReadAll(req.Body)

	body += fmt.Sprintf("Post body: %s\r\n", reqData)

	res.Write([]byte(body))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, mainPage)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
