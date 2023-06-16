package main

import (
	"database/sql"
	"flag"
	"log"
    "html/template"
	"net/http"
	"os"

	"github.com/rostamborn/snippetbox/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
    errorLog *log.Logger
    infoLog *log.Logger
    snippets *mysql.SnippetModel
    templateCache map[string]*template.Template
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
    mess := flag.String("message", "maaaaate", "show message")
    flag.Parse()

    infoLogger := log.New(os.Stdout, "INFO: ", log.LUTC | log.Ltime | log.Ldate)
    errLogger := log.New(os.Stderr, "ERROR\t", log.Lshortfile | log.LUTC | log.Ltime)

    db, err := openDB(*dsn) 
    if err != nil {
        errLogger.Fatal(err)
    }
    defer db.Close()

    templateChache, err := newTemplateCache("./ui/html/")
    if err != nil {
        errLogger.Fatal(err)
    }

    app := &application{
        errorLog: errLogger,
        infoLog: infoLogger,
        snippets: &mysql.SnippetModel{DB: db},
        templateCache: templateChache,
    }
    mux := app.routes()    

    infoLogger.Printf("starting server on address %s\n", *addr)
    infoLogger.Println("message: ", *mess)

    srv := http.Server{ Addr: *addr, ErrorLog: errLogger,
        Handler: mux,
    }
    err = srv.ListenAndServe()
    if err != nil {
        errLogger.Fatal(err)
    }
}
