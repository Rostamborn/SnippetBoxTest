package main

import (
    "log"
    "net/http"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", home)
    mux.HandleFunc("/snippet", showSnippet)
    mux.HandleFunc("/snippet/create", createSnippet)

    log.Println("starting server on port :8000")
    err := http.ListenAndServe(":8000", mux)
    if err != nil {
        log.Fatal(err)
    }
}
