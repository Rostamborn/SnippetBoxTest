package main

import (
    "net/http"
    "github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
    mux := chi.NewRouter()

    // middleware stack
    // mux.Use(app.recoverPanic, app.logRequest, secureHeaders)

    // standard middleware chain
    mux.Group(func(r chi.Router) {
        r.Use(app.recoverPanic)
        r.Use(app.logRequest)
        r.Use(secureHeaders)
        r.Use(app.session.Enable)
        r.Get("/", app.home)
        r.Get("/snippet/create", app.createSnippetForm)
        r.Post("/snippet/create", app.createSnippet)
        // apparently there is like a global key value thing in Chi
        // so that's why I can access id in handlers and get it's value
        r.Get("/snippet/{id}", app.showSnippet)
    })



    fileServer := http.FileServer(http.Dir("./ui/static/"))
    mux.Handle("/static/", http.StripPrefix("/static", fileServer))
    return mux
}
