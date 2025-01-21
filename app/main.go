package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/joho/godotenv"
	"github.com/syrlramadhan/pendaftaran-coc/config"
	"github.com/syrlramadhan/pendaftaran-coc/controller"
	"github.com/syrlramadhan/pendaftaran-coc/repository"
	"github.com/syrlramadhan/pendaftaran-coc/service"
)

func main() {
	errEnv := godotenv.Load()
	if errEnv != nil {
		panic(errEnv)
	}
	appPort := os.Getenv("APP_PORT")
	fmt.Println("listened and serve to port", appPort)

	sqlite, err := config.ConnectToDatabase()
	if err != nil {
		panic(err)
	}

	pendaftarRepository := repository.NewPendaftarRepository()
	pendaftarService := service.NewPendaftarService(pendaftarRepository, sqlite)
	pendaftarController := controller.NewPendaftarController(pendaftarService)

	router := httprouter.New()

	// add
	router.POST("/api/pendaftar/add", pendaftarController.CreatePendaftar)

	//login
	router.POST("/api/admin/login", pendaftarController.LoginAdmin)

	//get
	router.GET("/api/pendaftar/get", pendaftarController.ReadPendaftar)

	//serve file
	router.ServeFiles("/api/pendaftar/uploads/*filepath", http.Dir("uploads"))

	//download db
	router.ServeFiles("/api/pendaftar/database/*filepath", http.Dir("database"))

	handler := corsMiddleware(router)

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", appPort),
		Handler: handler,
	}

	errServe := server.ListenAndServe()
	if errServe != nil {
		panic(errServe)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
