package server

import (
	"github.com/williamnoble/kube-botany/pkg/plant"
	"time"
)

type PlantDTO struct {
	Id          string    `json:"Id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Age         string    `json:"age"`
	Growth      int64     `json:"growth"`
	GrowthStage string    `json:"stage"`
	WaterLevel  int       `json:"water"`
	WateredLast time.Time `json:"watered_last"`
	// Used when querying images
	Image     string `json:"image,omitempty"`
	DaysAlive int    `json:"days_alive,omitempty"`
}

func (s *Server) plantDTO(p *plant.Plant) PlantDTO {
	r := PlantDTO{
		Id:          p.Id,
		Name:        p.Name,
		Type:        string(p.Type),
		GrowthStage: p.GrowthStage.String(),
		Age:         time.Since(p.CreationTime).Round(time.Second).String(),
		WaterLevel:  p.WaterLevel,
		WateredLast: p.LastWatered,
		Image:       "",
		DaysAlive:   p.DaysAlive(),
	}

	if r.Id == "DefaultSunflower234" {
		r.WaterLevel = 10 // test unwatered plant
	}

	return r
}

type WaterResponse struct {
	Message string `json:"message"`
	Plant   PlantDTO
}

type WaterRequest struct {
	Id string `json:"id"` // NamespacedName
}
