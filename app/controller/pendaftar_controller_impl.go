package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
	"github.com/syrlramadhan/pendaftaran-coc/dto"
	"github.com/syrlramadhan/pendaftaran-coc/service"
	"github.com/syrlramadhan/pendaftaran-coc/util"
)

type PendaftarControllerImpl struct {
	PendaftarService service.PendaftarService
}

func NewPendaftarController(pendaftarService service.PendaftarService) PendaftarController {
	return &PendaftarControllerImpl{
		PendaftarService: pendaftarService,
	}
}

func (p *PendaftarControllerImpl) CreatePendaftar(writter http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pendaftarRequest := dto.PendaftarRequest{}
	err := request.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(writter, "Unable to parse form", http.StatusBadRequest)
		return
	}
	pendaftarRequest = dto.PendaftarRequest{
		NamaLengkap: request.FormValue("nama-lengkap"),
		Email: request.FormValue("email"),
		NoTelp: request.FormValue("no-telp"),
		Framework: request.FormValue("framework"),
	}
	file, handler, err := request.FormFile("bukti-transfer")
	if err != nil {
		http.Error(writter, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileName := fmt.Sprintf("%s.jpeg", pendaftarRequest.Email)
	handler.Filename = fileName
	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
	}

	filePath := filepath.Join(uploadDir, handler.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		panic(err)
	}

	pendaftarRequest.BuktiTransfer = handler.Filename

	responseDTO := p.PendaftarService.CreatePendaftar(request.Context(), pendaftarRequest)
	response := dto.ResponseList{
		Code:    http.StatusOK,
		Message: "success to add pendaftar",
		Data:    responseDTO,
	}

	// writter.Header().Set("Content-Type", "application/json")
	util.WriteToResponseBody(writter, response)
}

func (p *PendaftarControllerImpl) ReadPendaftar(writter http.ResponseWriter, request *http.Request, params httprouter.Params) {
	responseDTO := p.PendaftarService.ReadPendaftar(request.Context())
	response := dto.ResponseList{
		Code:    http.StatusOK,
		Message: "success to get data",
		Data:    responseDTO,
	}

	// writter.Header().Set("Content-Type", "application/json")
	util.WriteToResponseBody(writter, response)
}

func (p *PendaftarControllerImpl) LoginAdmin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var adminRequest dto.AdminRequest

	err := json.NewDecoder(r.Body).Decode(&adminRequest)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	token, err := p.PendaftarService.LoginAdmin(r.Context(), adminRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	response := dto.ResponseToken{
		Code:    http.StatusOK,
		Status:  "OK",
		Token:   token,
		Message: "token generate successfully",
	}

	w.Header().Set("Content-.Type", "application/json")
	util.WriteToResponseBody(w, response)
}