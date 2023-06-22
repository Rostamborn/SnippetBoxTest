package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rostamborn/snippetbox/pkg/models"
)

func secureHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        w.Header().Set("X-Frame-Options", "deny")
        next.ServeHTTP(w, r)
    })
}

func (app *application) logRequest(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.Host)
        next.ServeHTTP(w, r)
    })
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                w.Header().Set("Connection", "close")
                app.serveError(w, fmt.Errorf("%s", err))
            }
        }()
        next.ServeHTTP(w, r)
    })
}

func (app *application) requireAuthenticatedUser(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if app.authenticatedUser(r) == nil {
            http.Redirect(w, r, "/user/login", http.StatusFound)
            return
        }
    })
}

func (app *application) authenticate(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        exists := app.session.Exists(r, "userID")
        if !exists {
            next.ServeHTTP(w, r)
            return
        }

        user, err := app.users.Get(app.session.GetInt(r, "userID"))
        if err == models.ErrNorecord {
            app.session.Remove(r, "userID")
            next.ServeHTTP(w, r)
            return
        } else if err != nil {
            app.serveError(w, err)
            return
        }

        ctx := context.WithValue(r.Context(), contextKeyUser, user)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
