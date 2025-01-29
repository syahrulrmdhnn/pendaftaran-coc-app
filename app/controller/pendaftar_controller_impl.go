package controller

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/julienschmidt/httprouter"
	"github.com/syrlramadhan/pendaftaran-coc/app/model"
	"github.com/syrlramadhan/pendaftaran-coc/app/service"
)

type pendaftarControllerImpl struct {
	PendaftarService service.PendaftarService
}

func NewPendaftarController(pendaftarService service.PendaftarService) PendaftarController {
	return &pendaftarControllerImpl{
		PendaftarService: pendaftarService,
	}
}

// RenderTemplate implements PendaftarController.
func (p *pendaftarControllerImpl) RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	tmplPath := "templates/" + tmpl
	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
	}
	t, err := template.New(tmpl).Funcs(funcMap).ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Template tidak ditemukan: "+err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

// CreatePendaftar implements PendaftarController.
func (p *pendaftarControllerImpl) CreatePendaftar(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	err := request.ParseMultipartForm(10 << 20)
	if err != nil {
		// http.Error(writer, "unable to parse form", http.StatusBadRequest)
		http.Redirect(writer, request, "/form", http.StatusSeeOther)
		return
	}
	namaLengkap := request.FormValue("nama-lengkap")
	email := request.FormValue("email")
	noTelp := request.FormValue("no-telp")
	framework := request.FormValue("framework")
	file, header, err := request.FormFile("buktitf")
	if err != nil {
		http.Redirect(writer, request, "/form", http.StatusSeeOther)
		return
	}
	defer file.Close()

	fileName := fmt.Sprintf("%s.jpeg", email)
	header.Filename = fileName
	uploadDir := "assets/buktitf/"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
	}

	filePath := filepath.Join(uploadDir, header.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		http.Redirect(writer, request, "/form", http.StatusSeeOther)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Redirect(writer, request, "/form", http.StatusSeeOther)
		return
	}

	buktiTransfer := header.Filename

	pendaftar := model.Pendaftar{
		NamaLengkap:   namaLengkap,
		Email:         email,
		NoTelp:        noTelp,
		BuktiTransfer: buktiTransfer,
		Framework:     framework,
	}

	p.PendaftarService.CreatePendaftar(request.Context(), pendaftar)
	http.Redirect(writer, request, "/", http.StatusSeeOther)
}

// ReadPendaftar implements PendaftarController.
func (p *pendaftarControllerImpl) ReadPendaftar(request *http.Request, params httprouter.Params) ([]model.Pendaftar, error) {
	pendaftar, err := p.PendaftarService.ReadPendaftar(request.Context())
	if err != nil {
		panic(err)
	}

	return pendaftar, nil
}

// LoginAdmin implements PendaftarController.
func (p *pendaftarControllerImpl) LoginAdmin(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	user := request.FormValue("username")
	pass := request.FormValue("password")

	token, err := p.PendaftarService.LoginAdmin(request.Context(), user, pass)
	if err != nil {
		http.Redirect(writer, request, "/login", http.StatusSeeOther)
		return
	}

	http.SetCookie(writer, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true, // Agar lebih aman
		MaxAge: 600,
	})

	// data := map[string]string{"Token": token}
	http.Redirect(writer, request, "/pendaftar", http.StatusSeeOther)
	// p.RenderTemplate(writer, "pendaftar.html", data)
}
