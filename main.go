package main

import (
	"fmt"
	"net/http"
	"os"

	// "github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/syrlramadhan/pendaftaran-coc/app"
	"github.com/syrlramadhan/pendaftaran-coc/app/config"
)

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	panic(err)
	// }
	port := os.Getenv("APP_PORT")
	fmt.Println("runnig on port", port)

	sqlite, err := config.ConnectToDatabase()
	if err != nil {
		panic(err)
	}

	router := httprouter.New()

	router.ServeFiles("/assets/*filepath", http.Dir("assets"))
	router.ServeFiles("/database/*filepath", http.Dir("database"))

	handler := app.Routes(router, sqlite)

	server := http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	errServe := server.ListenAndServe()
	if errServe != nil {
		panic(errServe)
	}
}
