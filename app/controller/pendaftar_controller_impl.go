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
	"github.com/syrlramadhan/pendaftaran-coc/exception"
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

func writeJSONError(w http.ResponseWriter, code int, message string) {
	response := dto.ResponseError{
		Code:    code,
		Message: message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

func (p *PendaftarControllerImpl) CreatePendaftar(writter http.ResponseWriter, request *http.Request, params httprouter.Params) {
	defer func() {
		if r := recover(); r != nil {
			var response dto.ResponseError
			switch err := r.(type) {
			case exception.BadRequestHandler:
				writter.WriteHeader(http.StatusBadRequest)
				response = dto.ResponseError{
					Code:    http.StatusBadRequest,
					Message: err.Error(),
				}
			default:
				writter.WriteHeader(http.StatusInternalServerError)
				response = dto.ResponseError{
					Code:    http.StatusInternalServerError,
					Message: "internal server error",
				}
			}
			util.WriteToResponseBody(writter, response)
		}
	}()

	pendaftarRequest := dto.PendaftarRequest{}
	err := request.ParseMultipartForm(10 << 20)
	if err != nil {
		writeJSONError(writter, http.StatusBadRequest, "unable to parse form")
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
		writeJSONError(writter, http.StatusBadRequest, "failed to read file")
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
		Message: "registration successful",
		Data:    responseDTO,
	}

	util.WriteToResponseBody(writter, response)
}

func (p *PendaftarControllerImpl) ReadPendaftar(writter http.ResponseWriter, request *http.Request, params httprouter.Params) {
	responseDTO := p.PendaftarService.ReadPendaftar(request.Context())
	response := dto.ResponseList{
		Code:    http.StatusOK,
		Message: "success to get data",
		Data:    responseDTO,
	}

	util.WriteToResponseBody(writter, response)
}

func (p *PendaftarControllerImpl) LoginAdmin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var adminRequest dto.AdminRequest

	err := json.NewDecoder(r.Body).Decode(&adminRequest)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid input")
		return
	}

	token, err := p.PendaftarService.LoginAdmin(r.Context(), adminRequest)
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, err.Error())
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