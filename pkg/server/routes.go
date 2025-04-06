package server

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (s *Server) Routes() http.Handler {
	r := chi.NewRouter()

	// handle static assets
	r.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("cmd/api/static"))))

	// handle API Endpoints
	r.Route("/api/plants", func(r chi.Router) {
		r.Get("/", s.HandleIndex)
		r.Get("/{id}", s.HandlePlant)
		r.Delete("/{id}", s.HandlePlantDelete)
	})

	// handle Web
	r.HandleFunc("GET /", s.HandleCards)
	r.HandleFunc("POST /water", s.HandleWater)
	r.HandleFunc("GET /{id}", s.HandlePlantDetail)

	//r.HandleFunc("POST /plant", s.handleNewPlant)
	//r.HandleFunc("GET /ascii", s.handleASCII)

	return r
}
