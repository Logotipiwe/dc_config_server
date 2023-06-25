package main

import "github.com/google/uuid"

type Namespace struct {
	Id   string
	Name string
}

type Service struct {
	Id   string
	Name string
}

type Property struct {
	Id          string
	Name        string
	Value       string
	NamespaceId string
	ServiceId   string
	Active      bool
}

func NewProperty(name, value, namespaceId, serviceId string) Property {
	return Property{uuid.New().String(), name, value, namespaceId, serviceId, true}
}

func NewService(name string) Service {
	return Service{uuid.New().String(), name}
}
func CreateNamespace(name string) Namespace {
	return Namespace{uuid.New().String(), name}
}
