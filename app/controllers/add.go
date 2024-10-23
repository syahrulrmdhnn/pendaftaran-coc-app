package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/syrlramadhan/pendaftaran-coc/config"
	"github.com/syrlramadhan/pendaftaran-coc/models"
)

func contains(str, substr string) bool {
    return strings.Contains(str, substr)
}

func AddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(50 << 20) // Limit 50MB for file size
	if err != nil {
		http.Error(w, `{"message":"Error parsing form data"}`, http.StatusBadRequest)
		return
	}

	namaLengkap := r.FormValue("nama_lengkap")
	email := r.FormValue("email")
	telepon := r.FormValue("telepon")
	framework := r.FormValue("framework")

	// Validasi input
	if len(namaLengkap) < 2 {
		http.Error(w, `{"message":"Nama terlalu pendek"}`, http.StatusBadRequest)
		return
	}

	if len(email) < 2 || !contains(email, "@") {
		http.Error(w, `{"message":"Email tidak valid"}`, http.StatusBadRequest)
		return
	}

	if len(telepon) < 2 {
		http.Error(w, `{"message":"Nomor telepon terlalu pendek"}`, http.StatusBadRequest)
		return
	}

	// Cek apakah email sudah terdaftar
	var existingPendaftar models.Pendaftar
	if err := config.DB.Where("email = ?", email).First(&existingPendaftar).Error; err == nil {
		http.Error(w, `{"message":"Email sudah terdaftar"}`, http.StatusConflict)
		return
	}

	// Handle file upload
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, `{"message":"Error mengunggah file"}`, http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Simpan file bukti transfer
	filePath := fmt.Sprintf("static/bukti_tf%s.jpg", email)
	outFile, err := os.Create(filePath)
	if err != nil {
		http.Error(w, `{"message":"Gagal membuat file"}`, http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, `{"message":"Gagal menyimpan file"}`, http.StatusInternalServerError)
		return
	}

	// Insert new record
	pendaftar := models.Pendaftar{
		NamaLengkap:   namaLengkap,
		Email:         email,
		NoTelp:        telepon,
		BuktiTransfer: fmt.Sprintf("bukti_tf%s.jpg", email),
		Framework: framework,
	}
	if err := config.DB.Create(&pendaftar).Error; err != nil {
		http.Error(w, `{"message":"Gagal menyimpan data ke database"}`, http.StatusInternalServerError)
		return
	}

	w.Write([]byte(`{"message":"success"}`))
}