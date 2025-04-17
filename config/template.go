package config

import (
	"html/template"
)

var GLOBAL_TEMPLATE *template.Template

func InitTemplate() {
	temp, err := template.ParseGlob("./static/templates/*.html")
	GLOBAL_TEMPLATE = template.Must(temp, err)
}
