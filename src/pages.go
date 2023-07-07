package main

import (
	"bytes"
	"html/template"
)

func getAdminPage() (string, error) {
	tmpl := template.Must(template.ParseFiles("src/templates/index.gohtml"))

	props, _ := GetAllProps()
	namespaces, _ := GetAllNamespaces()
	services, _ := GetAllServices()
	view, err := CreateIndexView(props, namespaces, services)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, view); err != nil {
		return "", err
	}

	return tpl.String(), nil
}
