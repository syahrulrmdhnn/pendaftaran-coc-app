package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/syrlramadhan/pendaftaran-coc/config"
	"github.com/syrlramadhan/pendaftaran-coc/models"
)

func OrangHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["pendaftar"]

	var pendaftar models.Pendaftar
	if err := config.DB.Where("email = ?", email).First(&pendaftar).Error; err != nil {
		response := fmt.Sprintf(`{"nama_lengkap":"%s","email":"Tidak ada","telepon":"Tidak ada","bukti_tf":"Tidak ada","status":"Tidak Terdaftar"}`, email)
		w.Write([]byte(response))
		return
	}

	response := fmt.Sprintf(`{"nama_lengkap":"%s","email":"%s","telepon":"%s","bukti_tf":"%s","status":"Berhasil Mendaftar"}`,
		pendaftar.NamaLengkap, pendaftar.Email, pendaftar.NoTelp, pendaftar.BuktiTransfer)
	w.Write([]byte(response))
}