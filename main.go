package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/syrlramadhan/pendaftaran-coc/app"
	"github.com/syrlramadhan/pendaftaran-coc/app/config"
)

func main() {
	fmt.Println("Success to connect")

	sqlite, err := config.ConnectToDatabase()
	if err != nil {
		panic(err)
	}

	router := httprouter.New()

	router.ServeFiles("/assets/*filepath", http.Dir("assets"))
	router.ServeFiles("/database/*filepath", http.Dir("database"))

	handler := app.Routes(router, sqlite)

	server := http.Server{
		Addr: ":9000",
		Handler: handler,
	}

	errServe := server.ListenAndServe()
	if errServe != nil {
		panic(errServe)
	}
}