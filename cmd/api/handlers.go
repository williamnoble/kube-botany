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
		message = fmt.Sprintf("added %d units of water to %s (%s watered).", unitsAdded, plant.ID, plant.WaterLevelPercent())
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
