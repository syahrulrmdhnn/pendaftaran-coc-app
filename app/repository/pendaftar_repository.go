package repository

import (
	"context"
	"database/sql"

	"github.com/syrlramadhan/pendaftaran-coc/app/model"
)

type PendaftarRepository interface {
	CreatePendaftar(ctx context.Context, tx *sql.Tx, pendaftar model.Pendaftar) (model.Pendaftar, error)
	ReadPendaftar(ctx context.Context, tx *sql.Tx) []model.Pendaftar
	EmailExists(ctx context.Context, tx *sql.Tx, email string) bool
}