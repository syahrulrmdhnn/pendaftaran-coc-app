package service

import (
	"context"

	"github.com/syrlramadhan/pendaftaran-coc/dto"
)

type PendaftarService interface {
	CreatePendaftar(ctx context.Context, pendaftarRequest dto.PendaftarRequest) dto.PendaftarResponse
	ReadPendaftar(ctx context.Context) []dto.PendaftarResponse
	LoginAdmin(ctx context.Context, adminRequest dto.AdminRequest) (string, error)
}