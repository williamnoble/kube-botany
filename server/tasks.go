package server

import (
	"github.com/williamnoble/kube-botany/gen"
	"time"
)

func (s *Server) BackgroundTasks() {
	s.Logger.Info("starting background tasks")

	imgSvc := gen.NewMockImageGenerationService(s.staticDir, s.Logger)
	// run the task once on startup
	plants := s.store.ListAllPlants()
	err := imgSvc.ImageTask(plants)
	if err != nil {
		s.Logger.Error("error", err)
	}

	ticker := time.NewTimer(24 * time.Hour)
	defer ticker.Stop()
	for range ticker.C {
		err = imgSvc.ImageTask(plants)
		if err != nil {
			s.Logger.Error("error", err)
		}
	}
}
