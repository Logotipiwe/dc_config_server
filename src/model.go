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

func CreateProperty(name, value, namespaceId, serviceId string) Property {
	return NewProperty(uuid.NewString(), name, value, namespaceId, serviceId)
}

func NewProperty(id, name, value, namespaceId, serviceId string) Property {
	return Property{id, name, value, namespaceId, serviceId, true}
}

func CreateService(name string) Service {
	return Service{uuid.New().String(), name}
}

func CreateNamespace(name string) Namespace {
	return Namespace{uuid.New().String(), name}
}
