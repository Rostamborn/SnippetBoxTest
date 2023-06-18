package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rostamborn/snippetbox/pkg/forms"
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
    app.render(w, r, "create.page.tmpl", &templateData{
        Form: forms.New(nil),
    })
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
    // we explicitly call this to handle the potential errors
    // because we could've used r.PostFormValue that calles r.ParseForm automatically
    err := r.ParseForm()
    if err != nil {
        app.clientError(w, http.StatusBadRequest)
        return
    }
    //
    // title := r.PostForm.Get("title")
    // content := r.PostForm.Get("content")
    // expires := r.PostForm.Get("expires")
    //
    // errors := make(map[string]string)
    //
    // if strings.TrimSpace(title) == "" {
    //     errors["title"] = "This field cannot be blank"
    // } else if utf8.RuneCountInString(title) > 100 {
    //     errors["title"] = "This field is too long (maximum is 100 characters)"
    // }
    //
    // if strings.TrimSpace(content) == "" {
    //     errors["content"] = "This field cannot be blank"
    // }
    //
    // if strings.TrimSpace(expires) == "" {
    //     errors["expires"] = "This field cannot be blank"
    // } else if expires != "1" && expires != "7" && expires != "365" {
    //     errors["expires"] = "This field is invalid"
    // }
    //
    // if len(errors) > 0 {
    //     app.render(w, r, "create.page.tmpl", &templateData{
    //         FormErrors: errors,
    //         FormData: r.PostForm,
    //     })
    //     return
    // }
    form := forms.New(r.PostForm)
    form.Required("title", "content", "expires")
    form.MaxLength("title", 100)
    form.PermittedValues("expires", "365", "7", "1")

    if !form.Valid() {
        app.render(w, r, "create.page.tmpl", &templateData{Form: form})
    }

    id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
    if err != nil {
        app.serveError(w, err)
        return
    }

    // we redirect to get feedback on our insert request
    http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
