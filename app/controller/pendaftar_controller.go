package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/syrlramadhan/pendaftaran-coc/app/model"
)

type PendaftarController interface {
	RenderTemplate(writer http.ResponseWriter, tmpl string, data interface{})
	CreatePendaftar(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	ReadPendaftar(request *http.Request, params httprouter.Params) ([]model.Pendaftar, error)
	LoginAdmin(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}