package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/syrlramadhan/pendaftaran-coc/config"
	"github.com/syrlramadhan/pendaftaran-coc/models"
)

type Response struct {
	Message string                      `json:"message"`
	Data    map[string]models.Pendaftar `json:"data"`
}

func AmbilHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	nama := vars["nama"]
	kunci := vars["kunci"]

	n := os.Getenv("APP_NAMA")
	k := os.Getenv("APP_KUNCI")

	if nama != n || kunci != k {
		response := Response{
			Message: "Anda tidak memiliki akses!",
		}
		jsonData, err := json.MarshalIndent(response, "", "    ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
		return
	} else {
		var pendaftar []models.Pendaftar
		config.DB.Find(&pendaftar)
	
		data := make(map[string]models.Pendaftar)
		for _, p := range pendaftar {
			key := fmt.Sprintf("%d", p.ID)
			data[key] = models.Pendaftar{
				NamaLengkap:   p.NamaLengkap,
				Email:         p.Email,
				NoTelp:        p.NoTelp,
				BuktiTransfer: "https://pendaftaran-coc-api-production.up.railway.app/static/" + p.BuktiTransfer,
				Framework: p.Framework,
			}
		}
	
		response := Response{
			Message: "success",
			Data:    data,
		}
	
		jsonData, err := json.MarshalIndent(response, "", "    ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}

}
