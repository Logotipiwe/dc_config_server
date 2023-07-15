package main

import (
	"github.com/google/uuid"
	. "github.com/logotipiwe/dc_go_utils/src"
)

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

type CSPropertyDto DcPropertyDto

type NamespaceDto struct {
	Id   string
	Name string
}

type ServiceDto struct {
	Id   string
	Name string
}

type ImportExportAnswer struct {
	Namespaces []NamespaceDto  `json:"namespaces"`
	Services   []ServiceDto    `json:"services"`
	Props      []CSPropertyDto `json:"props"`
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

func (p Property) toDto() CSPropertyDto {
	return CSPropertyDto{
		Id:          p.Id,
		Name:        p.Name,
		Value:       p.Value,
		NamespaceId: p.NamespaceId,
		ServiceId:   p.ServiceId,
		Active:      p.Active,
	}
}

func (s *Service) toDto() ServiceDto {
	return ServiceDto{
		Id:   s.Id,
		Name: s.Name,
	}
}

func (s *Namespace) toDto() NamespaceDto {
	return NamespaceDto{
		Id:   s.Id,
		Name: s.Name,
	}
}

func csPropToModel(p CSPropertyDto) Property {
	return p.toModel()
}

func (p CSPropertyDto) toModel() Property {
	return Property{
		Id:          p.Id,
		Name:        p.Name,
		Value:       p.Value,
		NamespaceId: p.NamespaceId,
		ServiceId:   p.ServiceId,
		Active:      p.Active,
	}
}

func (s ServiceDto) toModel() Service {
	return Service{
		Id:   s.Id,
		Name: s.Name,
	}
}

func (n NamespaceDto) toModel() Namespace {
	return Namespace{
		Id:   n.Id,
		Name: n.Name,
	}
}
