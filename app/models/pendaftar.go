package models

// Pendaftar adalah model untuk tabel pendaftar
type Pendaftar struct {
	ID            uint   `gorm:"primaryKey"`
	NamaLengkap   string `gorm:"not null"`
	Email         string `gorm:"unique;not null"`
	NoTelp        string `gorm:"not null"`
	BuktiTransfer string
}