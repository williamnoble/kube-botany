package server

import "time"

func (s *Server) BackgroundTasks() {
	imgSvc := NewMockImageGenerationService(s.staticDir, s.Logger)
	// run the task once on startup
	err := imgSvc.imageTask(s.plants)
	if err != nil {
		s.Logger.Error("error", err)
	}

	ticker := time.NewTimer(24 * time.Hour)
	defer ticker.Stop()
	for range ticker.C {
		err = imgSvc.imageTask(s.plants)
		if err != nil {
			s.Logger.Error("error", err)
		}
	}
}
