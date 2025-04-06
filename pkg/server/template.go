package server

import (
	"html/template"
	"path/filepath"
)

func (s *Server) ParseTemplates() {
	s.templates["index"] = template.Must(template.ParseFiles(
		filepath.Join(s.TemplatesDir, "/layout.html"),
		filepath.Join(s.TemplatesDir, "/index.html")))

	s.templates["plant"] = template.Must(template.ParseFiles(
		filepath.Join(s.TemplatesDir, "/layout.html"),
		filepath.Join(s.TemplatesDir, "/plant.html")))
}
