package dto

type PendaftarRequest struct {
	NamaLengkap   string `json:"nama-lengkap"`
	Email         string `json:"email"`
	NoTelp        string `json:"no-telp"`
	BuktiTransfer string `json:"bukti-transfer"`
	Framework     string `json:"framework"`
}
	