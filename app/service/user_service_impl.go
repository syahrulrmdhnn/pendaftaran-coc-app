package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	// "github.com/joho/godotenv"
	"github.com/syrlramadhan/pendaftaran-coc/app/model"
	"github.com/syrlramadhan/pendaftaran-coc/app/repository"
)

var jwtKey = []byte("secret_key")

type PendaftarServiceImpl struct {
	PendaftarRepository repository.PendaftarRepository
	DB                  *sql.DB
	SecretKey           string
}

func NewPendaftarServiceImpl(pendaftarRepository repository.PendaftarRepository, db *sql.DB) PendaftarService {
	return &PendaftarServiceImpl{
		PendaftarRepository: pendaftarRepository,
		DB:                  db,
		SecretKey:           "secret_key",
	}
}

// CreateUser : Fungsi untuk melakukan validasi dan logika pada program.
// contohnya jika anda di suruh untuk melakukan validasi untuk pengecekan nomor hp yang tidak boleh sama di dalam table mst_user
func (pendaftarService PendaftarServiceImpl) CreatePendaftar(ctx context.Context, pendaftarModel model.Pendaftar) model.Pendaftar {
	tx, err := pendaftarService.DB.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Commit()
	emailExists := pendaftarService.PendaftarRepository.EmailExists(ctx, tx, pendaftarModel.Email)

	if emailExists {
		panic("Email sudah terdaftar")
	}

	uuidPendaftar := uuid.New().String()

	pendaftar := model.Pendaftar{
		Id:            uuidPendaftar,
		NamaLengkap:   pendaftarModel.NamaLengkap,
		Email:         pendaftarModel.Email,
		NoTelp:        pendaftarModel.NoTelp,
		Kampus:        pendaftarModel.Kampus,
		Alamat:        pendaftarModel.Alamat,
		BuktiTransfer: pendaftarModel.BuktiTransfer,
	}

	insertPendaftar, err := pendaftarService.PendaftarRepository.CreatePendaftar(ctx, tx, pendaftar)
	if err != nil {
		panic(err)
	}

	return insertPendaftar
}

func (pendaftarService PendaftarServiceImpl) ReadPendaftar(ctx context.Context) ([]model.Pendaftar, error) {
	tx, err := pendaftarService.DB.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Commit()
	users := pendaftarService.PendaftarRepository.ReadPendaftar(ctx, tx)

	return users, nil
}

type Claims struct {
	User string `json:"username"`
	jwt.StandardClaims
}

// GenerateJWT implements PendaftarService.
func (pendaftarService *PendaftarServiceImpl) GenerateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &Claims{
		User: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "admin-auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtKey)
}

// LoginAdmin implements PendaftarService.
func (pendaftarService *PendaftarServiceImpl) LoginAdmin(ctx context.Context, user string, pass string) (string, error) {
	// errEnv := godotenv.Load()
	// if errEnv != nil {
	// 	panic(errEnv)
	// }
	userAdmin := os.Getenv("USER_ADMIN")
	passAdmin := os.Getenv("PASS_ADMIN")
	tx, err := pendaftarService.DB.Begin()
	if err != nil {
		return "", fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Commit()
	if userAdmin == user && passAdmin == pass {
		fmt.Println("Login berhasil!")
	} else {
		return "", fmt.Errorf("invalid username or password")
	}

	token, err := pendaftarService.GenerateJWT(user)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return token, nil
}

func (pendaftarService *PendaftarServiceImpl) ValidateToken(ctx context.Context, tokenString string) (bool, error) {
	if tokenString == "" {
		return false, errors.New("token kosong")
	}

	// Parse token menggunakan secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metode signing tidak valid")
		}
		return []byte(pendaftarService.SecretKey), nil
	})

	if err != nil || !token.Valid {
		return false, errors.New("token tidak valid")
	}

	return true, nil
}
