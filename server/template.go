package server

import (
	"github.com/williamnoble/kube-botany/static/templates"
	"html/template"
)

// ParseTemplates parses HTML templates and stores them in the httpServer's templates map
// It parses the "index" template from layout.html and index.html
// It parses the "plant" template from layout.html and plant.html
func (s *Server) ParseTemplates() {
	//s.templates["index"] = template.Must(template.ParseFiles(
	//	filepath.Join(s.templatesDir, "/layout.html"),
	//	filepath.Join(s.templatesDir, "/index.html")))
	//
	//s.templates["plant"] = template.Must(template.ParseFiles(
	//	filepath.Join(s.templatesDir, "/layout.html"),
	//	filepath.Join(s.templatesDir, "/plant.html")))

	s.templates["index"] = template.Must(template.ParseFS(templates.HtmlTemplates,
		"layout.html",
		"index.html"))

	// Parse plant template (layout.html + plant.html)
	s.templates["plant"] = template.Must(template.ParseFS(templates.HtmlTemplates,
		"layout.html",
		"plant.html"))
}

func getTemplate() (*template.Template, error) {
	return template.ParseFS(templates.HtmlTemplates, "index.html")
}
