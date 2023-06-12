package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
    addr := flag.String("addr", ":8000", "HTTP network address")
    mess := flag.String("message", "maaaaate", "show message")
    flag.Parse()

    infoLogger := log.New(os.Stdout, "INFO: ", log.LUTC | log.Ltime | log.Ldate)
    errLogger := log.New(os.Stderr, "ERROR\t", log.Lshortfile | log.LUTC | log.Ltime)

    app := &application{
        errorLog: errLogger,
        infoLog: infoLogger,
    }
    mux := app.routes()    

    infoLogger.Printf("starting server on address %s\n", *addr)
    infoLogger.Println("message: ", *mess)

    srv := http.Server{
        Addr: *addr,
        ErrorLog: errLogger,
        Handler: mux,
    }
    err := srv.ListenAndServe()
    if err != nil {
        errLogger.Fatal(err)
    }
}
