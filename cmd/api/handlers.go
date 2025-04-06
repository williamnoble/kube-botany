package main

import (
	"fmt"
	"net/http"
	"time"
)

func (s *Server) HandleIndex(w http.ResponseWriter, r *http.Request) {
	var response []PlantResponse
	for _, p := range s.plants {
		p.Update(time.Now())
		plant := s.plantResponse(p)
		response = append(response, plant)
	}

	err := s.encode(w, r, http.StatusOK, response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (s *Server) HandleWater(w http.ResponseWriter, r *http.Request) {
	var waterReq WaterRequest
	err := s.decode(r, &waterReq)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	plant, err := s.plantByID(waterReq.Id)
	if err != nil {
		http.Error(w, "404", http.StatusBadRequest)
		return
	}

	message := "plant is fully watered and cannot be watered anymore."
	unitsAdded := plant.Water(time.Now())
	if unitsAdded > 0 {
		message = fmt.Sprintf("added %d units of water to %s (%s watered).", unitsAdded, plant.Id, plant.WaterLevelPercent())
	}

	response := WaterResponse{
		Message: message,
		Plant:   s.plantResponse(plant),
	}

	err = s.encode(w, r, http.StatusOK, response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (s *Server) HandleCards(w http.ResponseWriter, r *http.Request) {
	type NewPlant struct {
		ID         string
		Name       string
		Image      string
		Info       string
		WaterLevel int
	}

	data := []NewPlant{
		{ID: "1", Name: "my-bonsai-tree", Image: "/static/2025-04-06-bonsai123.png", Info: "bonsai", WaterLevel: 100},
		{ID: "2", Name: "little-sunflower", Image: "/static/2025-04-06-sunflower123.png", Info: "sunflower", WaterLevel: 60},
	}

	// Make sure you're executing the layout template
	err := s.templates["index"].ExecuteTemplate(w, "layout.html", data)
	if err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		s.logger.Error("template error", "error", err)
	}
}

// HandlePlantDetail Adds a new handler for plant detail pages
func (s *Server) HandlePlantDetail(w http.ResponseWriter, r *http.Request) {
	// Extract plant name from URL path
	plantName := r.PathValue("id")

	type NewPlant struct {
		ID         string
		Name       string
		Image      string
		Info       string
		WaterLevel int
		Day        int
	}

	// Find the plant with matching Id
	var selectedPlant *NewPlant
	for _, p := range []NewPlant{
		{ID: "bonsai", Name: "my-bonsai-tree", Image: "/static/2025-04-06-bonsai123.png", Info: "A miniature tree in a small container", Day: 30, WaterLevel: 100},
		{ID: "sunflower", Name: "little-sunflower", Image: "/static/2025-04-06-sunflower123.png", Info: "A tall plant with bright yellow flowers", Day: 30, WaterLevel: 60},
	} {
		if p.Name == plantName {
			selectedPlant = &p
			fmt.Println("found selected plant")
			break
		}
	}

	if selectedPlant == nil {
		http.Error(w, "Plant not found", http.StatusNotFound)
		return
	}

	// Execute only the plant.html template
	err := s.templates["plant"].ExecuteTemplate(w, "layout.html", selectedPlant)
	if err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		s.logger.Error("template error", "error", err)
	}
}
