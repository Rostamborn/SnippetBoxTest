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
        
        // apparently there is like a global key value thing in Chi
        // so that's why I can access id in handlers and get it's value
        r.Get("/snippet/{id}", app.showSnippet)
        r.Get("/user/signup", app.signupUserForm)
        r.Post("/user/signup", app.signupUser)
        r.Get("/user/login", app.loginUserForm)
        r.Post("/user/login", app.loginUser)
        r.Post("/user/logout", app.logoutUser)
    })

    mux.Group(func(r chi.Router) {
        r.Use(app.recoverPanic)
        r.Use(app.logRequest)
        r.Use(secureHeaders)
        r.Use(app.session.Enable)
        r.Use(app.requireAuthenticatedUser)
        r.Get("/snippet/create", app.createSnippetForm)
        r.Post("/snippet/create", app.createSnippet)
    })
        
    fileServer := http.FileServer(http.Dir("./ui/static/"))
    mux.Handle("/static/", http.StripPrefix("/static", fileServer))
    return mux
}
