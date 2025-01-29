package service

import (
	"context"

	"github.com/syrlramadhan/pendaftaran-coc/app/model"
)

type PendaftarService interface {
	CreatePendaftar(ctx context.Context, pendaftar model.Pendaftar) model.Pendaftar
	ReadPendaftar(ctx context.Context) ([]model.Pendaftar, error)
	GenerateJWT(username string) (string, error)
	LoginAdmin(ctx context.Context, user string, pass string) (string, error)
	ValidateToken(ctx context.Context, token string) (bool, error)
}
