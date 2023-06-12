package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type application struct {
    errorLog *log.Logger
    infoLog *log.Logger
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }

    files := []string{
        "./ui/html/home.page.tmpl",
        "./ui/html/base.layout.tmpl",
        "./ui/html/footer.partial.tmpl",
    }
    ts, err := template.ParseFiles(files...)
    if err != nil {
        app.errorLog.Println(err.Error())
        http.Error(w, "Internal Server Error", 500)
        return
    }

    err = ts.Execute(w, nil)
    if err != nil {
        app.errorLog.Println(err.Error())
        http.Error(w, "Internal Server Error", 500)
    }
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.URL.Query().Get("id"))
    if err != nil || id < 1 {
        app.errorLog.Println(err.Error())
        http.NotFound(w, r)
        return
    }
    fmt.Fprintf(w, "displaying item: %d", id)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        w.Header().Set("Allow", "POST")
        http.Error(w, "Method Not Allowed", 405)
        return
    }
    fmt.Fprint(w, "creating snippets here!")
}
