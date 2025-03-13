package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/syrlramadhan/pendaftaran-coc/app/model"
)

type PendaftarRepositoryImpl struct {
}

func NewPendaftarRepository() PendaftarRepository {
	return &PendaftarRepositoryImpl{}
}

func (p *PendaftarRepositoryImpl) CreatePendaftar(ctx context.Context, tx *sql.Tx, pendaftar model.Pendaftar) (model.Pendaftar, error) {
	query := `INSERT INTO pendaftars (id, nama_lengkap, email, no_telp, kampus, alamat, bukti_transfer) VALUES(?, ?, ?, ?, ?, ?, ?)`

	_, err := tx.ExecContext(ctx, query, pendaftar.Id, pendaftar.NamaLengkap, pendaftar.Email, pendaftar.NoTelp, pendaftar.Kampus, pendaftar.Alamat, pendaftar.BuktiTransfer)
	if err != nil {
		return pendaftar, fmt.Errorf("failed while entering data: %v", err)
	}

	return pendaftar, nil
}

func (p *PendaftarRepositoryImpl) ReadPendaftar(ctx context.Context, tx *sql.Tx) []model.Pendaftar {
	query := `SELECT id, nama_lengkap, email, no_telp, kampus, alamat, bukti_transfer FROM pendaftars`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		panic(err)
	}

	var pendaftars []model.Pendaftar
	for rows.Next() {
		pendaftar := model.Pendaftar{}
		err := rows.Scan(&pendaftar.Id, &pendaftar.NamaLengkap, &pendaftar.Email, &pendaftar.NoTelp, &pendaftar.Kampus, &pendaftar.Alamat, &pendaftar.BuktiTransfer)
		if err != nil {
			panic(err)
		}
		pendaftars = append(pendaftars, pendaftar)
	}

	return pendaftars
}

func (p *PendaftarRepositoryImpl) EmailExists(ctx context.Context, tx *sql.Tx, email string) bool {
	var exists bool
	query := `SELECT COUNT(1) FROM pendaftars WHERE email = ?`
	err := tx.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		panic(err)
	}

	return exists
}