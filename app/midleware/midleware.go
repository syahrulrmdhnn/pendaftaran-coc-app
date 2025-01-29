package midleware

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/syrlramadhan/pendaftaran-coc/app/service"
)

func AuthMiddleware(next httprouter.Handle, authService service.PendaftarService) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// Ambil token dari Header Authorization atau Cookie
		authHeader := r.Header.Get("Authorization")
		var token string

		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			token = cookie.Value
		}

		// Panggil ValidateToken dari authService
		isValid, err := authService.ValidateToken(r.Context(), token)
		if err != nil || !isValid {
			// http.Error(w, "Unauthorized: Token tidak valid", http.StatusUnauthorized)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Jika token valid, lanjutkan ke handler berikutnya
		next(w, r, p)
	}
}