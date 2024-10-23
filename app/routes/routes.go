package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/syrlramadhan/pendaftaran-coc/controllers"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Routes() {
	r := mux.NewRouter()

	r.HandleFunc("/api/add", controllers.AddHandler).Methods("POST")
	r.HandleFunc("/api/{pendaftar}", controllers.OrangHandler).Methods("GET")
	r.HandleFunc("/api/get/{nama}/{kunci}", controllers.AmbilHandler).Methods("GET")

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.Handle("/", corsMiddleware(r))
}