package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"github.com/rostamborn/snippetbox/pkg/models/mysql"
)

type contextKey string

var contextKeyUser = contextKey("user")

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *mysql.SnippetModel
	users         *mysql.UserModel
	templateCache map[string]*template.Template
	session       *sessions.Session
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	addr := flag.String("addr", ":8000", "HTTP network address")
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key for encryption")
	mess := flag.String("message", "maaaaate", "show message")
	flag.Parse()

	infoLogger := log.New(os.Stdout, "INFO: ", log.LUTC|log.Ltime|log.Ldate)
	errLogger := log.New(os.Stderr, "ERROR:\t", log.Lshortfile|log.LUTC|log.Ltime)

	db, err := openDB(*dsn)
	if err != nil {
		errLogger.Fatal(err)
	}
	defer db.Close()

	templateChache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errLogger.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 6 * time.Hour
	// session.Secure = true

	app := &application{
		errorLog:      errLogger,
		infoLog:       infoLogger,
		snippets:      &mysql.SnippetModel{DB: db},
		users:         &mysql.UserModel{DB: db},
		templateCache: templateChache,
		session:       session,
	}

	infoLogger.Printf("starting server on address %s\n", *addr)
	infoLogger.Println("message: ", *mess)

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errLogger,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err = srv.ListenAndServe()
	if err != nil {
		errLogger.Fatal(err)
	}
}
