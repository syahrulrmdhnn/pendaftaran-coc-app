package dto

type PendaftarResponse struct {
	Id            string `json:"id"`
	NamaLengkap   string `json:"nama-lengkap"`
	Email         string `json:"email"`
	NoTelp        string `json:"no-telp"`
	BuktiTransfer string `json:"bukti-transfer"`
	Framework     string `json:"framework"`
}
