package api

import (
	"html/template"
)

var (
	tempProcessors = template.Must(template.ParseGlob("templates/*.html"))
)
