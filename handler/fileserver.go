package handler

import (
	"kerrigan-chart/config"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

func FileServer(router *chi.Mux) {
	root := config.FrontendDir
	fs := http.FileServer(http.Dir(root))
	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(root + r.RequestURI); os.IsNotExist(err) {
			http.StripPrefix(r.RequestURI, fs).ServeHTTP(w, r)
		} else {
			fs.ServeHTTP(w, r)
		}
	})
}
