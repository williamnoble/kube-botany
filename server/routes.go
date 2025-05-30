package server

import (
	chi "github.com/go-chi/chi/v5"
	"net/http"
)

// Routes sets up the HTTP routes for the httpServer
// It defines routes for static assets, API endpoints, and web pages
func (s *Server) Routes() http.Handler {
	r := chi.NewRouter()

	// handle static assets
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// handle API Endpoints
	r.Route("/api/plants", func(r chi.Router) {
		r.Get("/", s.HandleListPlants)         // GET /api/plants - List all plants
		r.Get("/{id}", s.HandleGetPlant)       // GET /api/plants/{id} - Get a specific plant
		r.Delete("/{id}", s.HandlePlantDelete) // DELETE /api/plants/{id} - Delete a plant
		r.Post("/water/{id}", s.HandleWaterPlant)
		r.Post("/", s.HandleCreatePlant) // POST /api/plants - Create a plant
	})

	// handle Web
	r.HandleFunc("GET /", s.HandleRenderHomePage)  // GET / - Render home page with all plants
	r.HandleFunc("GET /{id}", s.HandlePlantDetail) // GET /{id} - Render plant detail page

	// Commented out routes for future implementation
	//r.HandleFunc("POST /plant", s.handleNewPlant)
	//r.HandleFunc("GET /ascii", s.handleASCII)

	return r
}
