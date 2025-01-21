package service

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/syrlramadhan/pendaftaran-coc/dto"
	"github.com/syrlramadhan/pendaftaran-coc/model"
	"github.com/syrlramadhan/pendaftaran-coc/repository"
	"github.com/syrlramadhan/pendaftaran-coc/util"
)

var jwtKey = []byte("secret_key")

type PendaftarServiceImpl struct {
	PendaftarRepository repository.PendaftarRepository
	DB                  *sql.DB
}

func NewPendaftarService(pendaftarRepo repository.PendaftarRepository, db *sql.DB) PendaftarService {
	return &PendaftarServiceImpl{
		PendaftarRepository: pendaftarRepo,
		DB:                  db,
	}
}

func (p *PendaftarServiceImpl) CreatePendaftar(ctx context.Context, pendaftarRequest dto.PendaftarRequest) dto.PendaftarResponse {
	tx, err := p.DB.Begin()
	if err != nil {
		panic(err)
	}
	defer util.CommitOrRollBack(tx)

	pendaftar := model.Pendaftar{
		Id:            uuid.New().String(),
		NamaLengkap:   pendaftarRequest.NamaLengkap,
		Email:         pendaftarRequest.Email,
		NoTelp:        pendaftarRequest.NoTelp,
		BuktiTransfer: pendaftarRequest.BuktiTransfer,
		Framework:     pendaftarRequest.Framework,
	}
	createPendaftar, err := p.PendaftarRepository.CreatePendaftar(ctx, tx, pendaftar)
	if err != nil {
		panic(err)
	}

	return util.ToPendaftarResponse(createPendaftar)
}

func (p *PendaftarServiceImpl) ReadPendaftar(ctx context.Context) []dto.PendaftarResponse {
	tx, err := p.DB.Begin()
	if err != nil {
		panic(err)
	}
	pendaftar := p.PendaftarRepository.ReadPendaftar(ctx, tx)

	return util.ToPendaftarListResponse(pendaftar)
}

type Claims struct {
	User string `json:"username"`
	jwt.StandardClaims
}

func (p *PendaftarServiceImpl) GenerateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &Claims{
		User: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "go-auth-example",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtKey)
}

func (p *PendaftarServiceImpl) LoginAdmin(ctx context.Context, adminRequest dto.AdminRequest) (string, error) {
	errEnv := godotenv.Load()
	if errEnv != nil {
		panic(errEnv)
	}
	userAdmin := os.Getenv("USER_ADMIN")
	passAdmin := os.Getenv("PASS_ADMIN")
	tx, err := p.DB.Begin()
	if err != nil {
		return "", fmt.Errorf("failed to start transaction: %v", err)
	}
	defer util.CommitOrRollBack(tx)
	if userAdmin == adminRequest.User && passAdmin == adminRequest.Pass {
		fmt.Println("Login berhasil!")
	} else {
		return "", fmt.Errorf("invalid email or password")
	}

	token, err := p.GenerateJWT(adminRequest.User)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return token, nil
}
