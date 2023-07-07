package main

import (
	"encoding/json"
	. "github.com/logotipiwe/dc_go_utils/src"
)

type ServiceView struct {
	Id   string
	Name string
}

func (service *Service) toView() ServiceView {
	return ServiceView{service.Id, service.Name}
}

type NamespaceView struct {
	Id   string
	Name string
}

func (namespace *Namespace) toView() NamespaceView {
	return NamespaceView{namespace.Id, namespace.Name}
}

type PropertyView struct {
	Id            string
	ServiceId     string
	ServiceName   string
	NamespaceId   string
	NamespaceName string
	Name          string
	Value         string
	Active        bool
}

func (p Property) toView(namespace Namespace, service Service) PropertyView {
	return PropertyView{
		p.Id,
		p.ServiceId,
		service.Name,
		p.NamespaceId,
		namespace.Name,
		p.Name,
		p.Value,
		p.Active,
	}
}

type IndexView struct {
	Services   []ServiceView
	Namespaces []NamespaceView
	Properties []PropertyView
	PropsJson  string
}

func CreateIndexView(props []Property, namespaces []Namespace, services []Service) (IndexView, error) {
	nMap := toMap(namespaces, func(val Namespace) string { return val.Id })
	sMap := toMap(services, func(s Service) string { return s.Id })
	propViews := Map(props, func(p Property) PropertyView {
		n := nMap[p.NamespaceId]
		s := sMap[p.ServiceId]
		return p.toView(n, s)
	})
	propDtos := Map(props, func(p Property) CSPropertyDto {
		return p.toDto()
	})
	propsJson, err := json.Marshal(propDtos)
	if err != nil {
		return IndexView{}, err
	}
	return IndexView{
		Services: Map(services, func(s Service) ServiceView {
			return s.toView()
		}),
		Namespaces: Map(namespaces, func(n Namespace) NamespaceView {
			return n.toView()
		}),
		Properties: propViews,
		PropsJson:  string(propsJson),
	}, nil
}
