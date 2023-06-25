package main

import (
	"bytes"
	"html/template"
)

func getAdminPage() string {
	tmpl := template.Must(template.ParseFiles("src/templates/index.gohtml"))

	props, _ := GetAllProps()
	namespaces, _ := GetAllNamespaces()
	services, _ := GetAllServices()
	view := CreateIndexView(props, namespaces, services)

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, view); err != nil {
		return err.Error()
	}

	return tpl.String()
}
