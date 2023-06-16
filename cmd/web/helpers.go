package main

import (
	"fmt"
    "time"
    "bytes"
	"net/http"
	"runtime/debug"
)

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
    if td == nil {
        td = &templateData{}
    }

    td.CurrentYear = time.Now().Year()
    return td
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
    ts, ok := app.templateCache[name]
    if !ok {
        app.serveError(w, fmt.Errorf("the template %s does not exist", name))
        return
    }

    buf := new(bytes.Buffer)

    err := ts.Execute(buf, app.addDefaultData(td, r))
    if err != nil {
        app.serveError(w, err)
        return
    }

    buf.WriteTo(w)
}

func (app *application) serveError(w http.ResponseWriter, err error) {
    trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
    // we do this to indicate the origin of the error instead of this file
    app.errorLog.Output(2, trace)

    http.Error(
        w,
        http.StatusText(http.StatusInternalServerError),
        http.StatusInternalServerError,
    )
}

func (app *application) clientError(w http.ResponseWriter, status int) {
    http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
    app.clientError(w, http.StatusNotFound)
}
