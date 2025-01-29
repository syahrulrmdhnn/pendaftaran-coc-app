package main

import (
	"fmt"
	"net/http"
	"os"

	// "github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/syrlramadhan/pendaftaran-coc/app"
	"github.com/syrlramadhan/pendaftaran-coc/app/config"
	"github.com/syrlramadhan/pendaftaran-coc/app/midleware"
	"github.com/syrlramadhan/pendaftaran-coc/app/repository"
	"github.com/syrlramadhan/pendaftaran-coc/app/service"
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

	pendaftarRepository := repository.NewPendaftarRepository()
	pendaftarService := service.NewPendaftarServiceImpl(pendaftarRepository, sqlite)

	router := httprouter.New()

	router.ServeFiles("/assets/*filepath", http.Dir("assets"))
	router.GET("/database/*filepath", midleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		http.StripPrefix("/database/", http.FileServer(http.Dir("database"))).ServeHTTP(w, r)
	}, pendaftarService))

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
