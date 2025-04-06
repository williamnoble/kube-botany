package server

import (
	"html/template"
	"path/filepath"
)

// ParseTemplates parses HTML templates and stores them in the server's templates map
// It parses the "index" template from layout.html and index.html
// It parses the "plant" template from layout.html and plant.html
func (s *Server) ParseTemplates() {
	s.templates["index"] = template.Must(template.ParseFiles(
		filepath.Join(s.TemplatesDir, "/layout.html"),
		filepath.Join(s.TemplatesDir, "/index.html")))

	s.templates["plant"] = template.Must(template.ParseFiles(
		filepath.Join(s.TemplatesDir, "/layout.html"),
		filepath.Join(s.TemplatesDir, "/plant.html")))
}
