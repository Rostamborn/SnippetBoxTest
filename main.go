package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }
    w.Write([]byte("hello mate"))
}

func showSnippet(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("we show snippets here"))
}

func createSnippet(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        // we have to set header before we call w.WriteHeader or w.Write
        // for it to have effect
        w.Header().Set("Allow", "POST")
        http.Error(w, "Method Not Allowed", 405)
        return
    }
    w.Write([]byte("we create sippets here"))
}

func main() {
    mux := http.NewServeMux()

    mux.HandleFunc("/", home)
    mux.HandleFunc("/snippet", showSnippet)
    mux.HandleFunc("/snippet/create", createSnippet)
    
    log.Println("Starting server on :4000")
    err := http.ListenAndServe(":4000", mux)
    if err != nil {
        log.Fatal(err)
    }
}
