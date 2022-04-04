package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sarmerer/go-crypto-dashboard/api/chartjs"
	"github.com/sarmerer/go-crypto-dashboard/config"
	"github.com/sarmerer/go-crypto-dashboard/repository"
)

func Serve(repo repository.Repository) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/api", func(r chi.Router) {
		r.Mount("/chartjs", chartjs.Route())
	})

	port := fmt.Sprintf(":%d", config.APIPort)
	log.Printf("listening on port %s", port)
	http.ListenAndServe(port, r)

	_ = repo
}
