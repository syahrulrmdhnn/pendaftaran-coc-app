package app

import (
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/syrlramadhan/pendaftaran-coc/app/controller"
	"github.com/syrlramadhan/pendaftaran-coc/app/midleware"
	"github.com/syrlramadhan/pendaftaran-coc/app/repository"
	"github.com/syrlramadhan/pendaftaran-coc/app/service"
)

func Routes(router *httprouter.Router, db *sql.DB) *httprouter.Router {

	pendaftarRepository := repository.NewPendaftarRepository()
	pendaftarService := service.NewPendaftarServiceImpl(pendaftarRepository, db)
	pendaftarController := controller.NewPendaftarController(pendaftarService)

	router.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		pendaftarController.RenderTemplate(w, "index.html", nil)
	})
	router.GET("/documentation", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		pendaftarController.RenderTemplate(w, "dokumentasi.html", nil)
	})
	router.GET("/form", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		pendaftarController.RenderTemplate(w, "formulir.html", nil)
	})
	router.POST("/form/add", pendaftarController.CreatePendaftar)

	// Gunakan middleware untuk proteksi halaman /pendaftar
	router.GET("/pendaftar", midleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		data, err := pendaftarController.ReadPendaftar(r, p)
		if err != nil {
			http.Error(w, "failed to read data", http.StatusInternalServerError)
			return
		}
		pendaftarController.RenderTemplate(w, "pendaftar.html", data)
	}, pendaftarService))

	router.GET("/login", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		pendaftarController.RenderTemplate(w, "login.html", nil)
	})
	router.POST("/login", pendaftarController.LoginAdmin)

	return router
}
