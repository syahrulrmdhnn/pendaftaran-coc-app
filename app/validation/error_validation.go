package validation

import (
	"context"
	"database/sql"

	"github.com/syrlramadhan/pendaftaran-coc/exception"
	"github.com/syrlramadhan/pendaftaran-coc/model"
	"github.com/syrlramadhan/pendaftaran-coc/repository"
)

func ValidateIfEmailExist(ctx context.Context, tx *sql.Tx, repo repository.PendaftarRepository, email string) {
	emailExist := repo.EmailExists(ctx, tx, email)
	if emailExist {
		panic(exception.NewBadRequestHandler("email already exists"))
	}
}

func ValidateIfNull(pendaftar model.Pendaftar) {
	if pendaftar.NamaLengkap == "" || pendaftar.NoTelp == "" || pendaftar.Email == "" || pendaftar.Framework == "" {
		panic(exception.NewBadRequestHandler("input cannot be empty"))
	}
}