# Website Pendaftaran Coconut Open Class

Deskripsi singkat tentang proyek Anda. Jelaskan tujuan dari proyek ini dan fitur-fitur utama yang disediakan.

## Prasyarat

Pastikan Anda sudah menginstal prasyarat berikut:

- [Golang](https://golang.org/dl/)
- [Git](https://git-scm.com/)

## Instalasi

1. Clone repository ini ke komputer Anda:

    ```bash
    git clone https://github.com/syahrulrmdhnn/pendaftaran-coc.git
    cd nama-proyek
    ```

2. Buat file `.env` di root folder proyek Anda. File ini digunakan untuk menyimpan konfigurasi lingkungan seperti database, API keys, dan sebagainya.

    Contoh isi `.env`:

    ```env
    APP_PORT = 5000
    APP_NAMA = coc02
    APP_KUNCI = coconut@013
    ```

3. Instal dependensi yang diperlukan:

    ```bash
    go mod tidy
    ```

## Menjalankan Aplikasi

Untuk menjalankan aplikasi, gunakan perintah berikut:

```bash
go run main.go
