package tool

import "html/template"

type Fill struct {
	Title string
	Core  template.HTML
	Style template.CSS
}
