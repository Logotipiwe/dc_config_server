package main

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
}

func CreateIndexView(props []Property, namespaces []Namespace, services []Service) IndexView {
	nMap := toMap(namespaces, func(val Namespace) string { return val.Id })
	sMap := toMap(services, func(s Service) string { return s.Id })
	return IndexView{
		Services: Map(services, func(s Service) ServiceView {
			return s.toView()
		}),
		Namespaces: Map(namespaces, func(n Namespace) NamespaceView {
			return n.toView()
		}),
		Properties: Map(props, func(p Property) PropertyView {
			n := nMap[p.NamespaceId]
			s := sMap[p.ServiceId]
			return p.toView(n, s)
		}),
	}
}
