package main

import (
    "net/http"
    "github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
    mux := chi.NewRouter()
    mux.Get("/", app.home)
    mux.Get("/snippet/create", app.createSnippetForm)
    mux.Post("/snippet/create", app.createSnippet)
    // apparently there is like a global key value thing in Chi
    // so that's why I can access id in handlers and get it's value
    mux.Get("/snippet/{id}", app.showSnippet)

    fileServer := http.FileServer(http.Dir("./ui/static/"))
    mux.Handle("/static/", http.StripPrefix("/static", fileServer))

    return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
