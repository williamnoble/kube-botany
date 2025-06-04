package server

import (
	"github.com/williamnoble/kube-botany/static/templates"
	"html/template"
)

// parseTemplate creates a new template from the embedded filesystem.
func (s *Server) parseTemplate(name string, files ...string) {
	s.templates[name] = template.Must(template.ParseFS(templates.HtmlTemplates, files...))
}

// ParseTemplates initialises all templates used by the server.
func (s *Server) ParseTemplates() {
	s.parseTemplate("index", "layout.html", "index.html")
	s.parseTemplate("plant", "layout.html", "plant.html")
}
