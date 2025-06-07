package templates

import (
	"embed"
)

//go:embed *.html
var HtmlTemplates embed.FS
