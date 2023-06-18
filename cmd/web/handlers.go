package main

import (
	"fmt"
	"net/http"
	"strconv"
    "github.com/go-chi/chi/v5"
	"github.com/rostamborn/snippetbox/pkg/models"
)



func (app *application) home(w http.ResponseWriter, r *http.Request) {
    s, err := app.snippets.Latest()
    if err != nil {
        app.serveError(w, err)
        return
    }

    data := &templateData{Snippets: s}
    app.render(w, r, "home.page.tmpl", data)
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
    // apparently there is like a global key value thing in Chi
    // so that's why I can use the key here that I defined in routes.go
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil || id < 1 {
        app.notFound(w)
        return
    }

    s, err := app.snippets.Get(id)
    if err == models.ErrNorecord {
        app.notFound(w)
        return
    } else if err != nil {
        app.serveError(w, err)
        return
    }

    // we wrap the snippet in a templateData struct to make it available in the template
    // also it allows us to use more dynamic data
    data := &templateData{Snippet: s}
    app.render(w, r, "show.page.tmpl", data)
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "we create a new snippet here")
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
    title := "bookshelf"
    content := "my bookshelf holds valueable books that are dear to me\nalas I'm dumb"
    expires := "8"

    id, err := app.snippets.Insert(title, content, expires)
    if err != nil {
        app.serveError(w, err)
        return
    }
    // we redirect to get feedback on our insert request
    // http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
    http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
