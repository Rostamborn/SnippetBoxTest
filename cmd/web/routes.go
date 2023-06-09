package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

    mux.Route("/", func(r chi.Router) {
        r.Use(app.recoverPanic)
		r.Use(app.logRequest)
		r.Use(secureHeaders)
		r.Use(app.session.Enable)
		r.Use(app.authenticate)
		r.Get("/", app.home)
		// apparently there is like a global key value thing in Chi
		// so that's why I can access id in handlers and get it's value
		r.Get("/snippet/{id}", app.showSnippet)
		r.Get("/user/signup", app.signupUserForm)
		r.Post("/user/signup", app.signupUser)
		r.Get("/user/login", app.loginUserForm)
		r.Post("/user/login", app.loginUser)
        r.Post("/user/logout", app.logoutUser)

        // Subroutes
        r.Route("/snippet/create", func(r chi.Router) {
            r.Use(app.requireAuthenticatedUser)
            r.Get("/", app.createSnippetForm)
            r.Post("/", app.createSnippet)
        })
    })

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return mux
}
