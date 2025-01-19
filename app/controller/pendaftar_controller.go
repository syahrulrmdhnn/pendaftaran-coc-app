package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type PendaftarController interface {
	CreatePendaftar(writter http.ResponseWriter, request *http.Request, params httprouter.Params)
	ReadPendaftar(writter http.ResponseWriter, request *http.Request, params httprouter.Params)
	LoginAdmin(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
}