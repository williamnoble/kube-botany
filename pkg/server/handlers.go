package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/williamnoble/kube-botany/pkg/plant"
	"net/http"
	"time"
)

// HandleListPlants returns a list of all plants as JSON
func (s *Server) HandleListPlants(w http.ResponseWriter, r *http.Request) {
	fmt.Println("sleep failed")
	var response []PlantDTO
	for _, p := range s.plants {
		p.Update(time.Now())
		plant := s.plantDTO(p)
		response = append(response, plant)
	}

	err := s.encode(w, r, http.StatusOK, response)
	if err != nil {
		http.Error(w, "Internal httpServer error", http.StatusInternalServerError)
	}
}

// HandleGetPlant returns a single plant by ID as JSON
func (s *Server) HandleGetPlant(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var found bool
	var response PlantDTO
	for _, p := range s.plants {
		if p.Id == id {
			found = true
			p.Update(time.Now())
			response = s.plantDTO(p)
		}
	}

	if found == false || id == "" {
		http.Error(w, "Plant not found, please check the id.", http.StatusNotFound)
		return
	}

	err := s.encode(w, r, http.StatusOK, response)
	if err != nil {
		http.Error(w, "Internal httpServer error", http.StatusInternalServerError)
	}
}

// HandlePlantDelete deletes a plant by ID
// It returns 204 No Content if successful, or 404 Not Found if the plant doesn't exist
func (s *Server) HandlePlantDelete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var found bool
	var index int
	for i, p := range s.plants {
		if p.Id == id {
			found = true
			p.Update(time.Now())
			index = i
			break
		}
	}

	if found == false || id == "" {
		http.Error(w, "Plant not found, please check the id.", http.StatusNotFound)
		return
	}

	s.plants = append(s.plants[:index], s.plants[index+1:]...)
	w.WriteHeader(http.StatusNoContent)
}

// HandleWaterPlant adds water to a plant
// It accepts a JSON request with a plant ID, adds water to the plant, and returns a response with a message and the updated plant
func (s *Server) HandleWaterPlant(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	plant, err := s.plantByID(id)
	if err != nil {
		http.Error(w, "Plant not found", http.StatusNotFound)
		return
	}

	message := "plant is fully watered and cannot be watered anymore."
	unitsAdded := plant.AddWater(time.Now())
	if unitsAdded > 0 {
		message = fmt.Sprintf("added %d units of water to %s (%d%% watered).", unitsAdded, plant.Id, plant.WaterLevel)
	}

	response := WaterResponse{
		Message: message,
		Plant:   s.plantDTO(plant),
	}

	err = s.encode(w, r, http.StatusOK, response)
	if err != nil {
		http.Error(w, "Internal httpServer error", http.StatusInternalServerError)
	}
}

// HandleRenderHomePage renders the home page with cards for all plants
// It converts each plant to a DTO, sets the image path, and renders the index.html template
func (s *Server) HandleRenderHomePage(w http.ResponseWriter, r *http.Request) {
	//time.Sleep(30 * time.Second)

	var data []PlantDTO
	for _, plant := range s.plants {
		dto := s.plantDTO(plant)
		dto.Image = fmt.Sprintf("/static/images/%s", plant.Image())
		fmt.Println("looking for image: ", plant.Image(), dto.Image)
		if dto.Name == "" {
			dto.Name = dto.Id
		}
		data = append(data, dto)
	}

	// Execute the index.html template with the layout
	err := s.templates["index"].ExecuteTemplate(w, "layout.html", data)
	if err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		s.Logger.Error("template error", "error", err)
	}
}

// HandlePlantDetail renders the plant detail page for a specific plant
// It extracts the plant ID from the URL path, finds the plant, and renders the plant.html template
func (s *Server) HandlePlantDetail(w http.ResponseWriter, r *http.Request) {
	// Extract plant ID from URL path
	plantId := r.PathValue("id")

	var dto PlantDTO
	var found bool
	for _, plant := range s.plants {
		if plant.Id == plantId {
			found = true
			dto = s.plantDTO(plant)
			dto.Image = fmt.Sprintf("/static/images/%s", plant.Image())
			if dto.Name == "" {
				dto.Name = dto.Id
			}
		}
	}

	if !found {
		http.Error(w, "Plant not found", http.StatusNotFound)
		return
	}

	// Execute the plant.html template with the layout
	err := s.templates["plant"].ExecuteTemplate(w, "layout.html", dto)
	if err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		s.Logger.Error("template error", "error", err)
	}
}

func (s *Server) HandleCreatePlant(w http.ResponseWriter, r *http.Request) {
	var dto PlantDTO
	err := s.decode(r, &dto)
	if err != nil {
		http.Error(w, "Error decoding request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	p := plant.NewPlant(
		dto.Id,
		dto.Name,
		plant.Sunflower, // TODO: Fix this typing
		dto.CreationTime,
		false,
	)

	s.plants = append(s.plants, p)
	w.WriteHeader(http.StatusCreated)

}

//
//func (s *Server) handleASCII(w http.ResponseWriter, r *http.Request) {
//	currentPlant.Update(time.Now())
//	currentPlant.GrowthStage = plant.Maturing
//	asciiArt := s.renderer.RenderFern(currentPlant)
//	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
//	w.WriteHeader(http.StatusOK)
//	w.Write([]byte(asciiArt))
//}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006 at 15:04")
}
