package server

import (
	"context"
	"github.com/williamnoble/kube-botany/pkg/gen"
	"time"
)

// BackgroundTasks sets up background tasks:
func (s *Server) BackgroundTasks(ctx context.Context) {
	s.Logger.With("component", "tasks").Info("starting background tasks")
	imgSvc := gen.NewMockImageGenerationService(s.staticDir, s.Logger)

	// Run the task once on startup
	if err := runImageTask(s, imgSvc); err != nil {
		s.Logger.With("component", "tasks").Error("error processing initial task", "error", err)
	}

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := runImageTask(s, imgSvc); err != nil {
				s.Logger.With("component", "tasks").Error("error processing scheduled task", "error", err)
			}
		case <-ctx.Done():
			s.Logger.With("component", "tasks").Info("stopping background tasks")
			return
		}
	}
}

// runImageTask runs the image generation task with the current list of plants
func runImageTask(s *Server, imgSvc *gen.ImageGenerationService) error {
	plants := s.store.ListAllPlants()
	return imgSvc.ImageTask(plants)
}
