package main

import (
	"fmt"
	// "html/template"
	"net/http"
	"strconv"

	"github.com/rostamborn/snippetbox/pkg/models"
)



func (app *application) home(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        app.notFound(w)
        return
    }

    s, err := app.snippets.Latest()
    if err != nil {
        app.serveError(w, err)
        return
    }

    data := &templateData{Snippets: s}
    app.render(w, r, "home.page.tmpl", data)
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.URL.Query().Get("id"))
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

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" { 
        w.Header().Set("Allow", "POST")
        app.clientError(w, http.StatusMethodNotAllowed)
        return
    }

    title := "bookshelf"
    content := "my bookshelf holds valueable books that are dear to me\nalas I'm dumb"
    expires := "8"

    id, err := app.snippets.Insert(title, content, expires)
    if err != nil {
        app.serveError(w, err)
        return
    }
    // we redirect to get feedback on our inser request
    http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
