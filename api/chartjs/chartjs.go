package chartjs

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Route() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Timeout(60))
	r.Use(middleware.StripSlashes)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ChartJS API"))
	})

	return r
}
