package datasets

import (
	"github.com/google/uuid"
	utils "github.com/logotipiwe/dc_go_utils/src"
)

var (
	emptyAns = []utils.DcPropertyDto{}
)

type NamespaceDto struct {
	Id   string
	Name string
}

type ServiceDto struct {
	Id   string
	Name string
}

type ImportExportAnswer struct {
	Namespaces []NamespaceDto        `json:"namespaces"`
	Services   []ServiceDto          `json:"services"`
	Props      []utils.DcPropertyDto `json:"props"`
}

type Dataset struct {
	Data       ImportExportAnswer
	Assertions []Assertion
	TestNames  []string
}

type Assertion struct {
	Service   string
	Namespace string
	Answer    []utils.DcPropertyDto
}

type Datasets struct {
}

func CreatePropFull(id, name, value, service, namespace string, active bool) utils.DcPropertyDto {
	return utils.DcPropertyDto{
		Id:          id,
		Name:        name,
		Value:       value,
		NamespaceId: namespace,
		ServiceId:   service,
		Active:      active,
	}
}

func CreateProp(service, namespace string, active bool) utils.DcPropertyDto {
	return utils.DcPropertyDto{
		Id:          uuid.NewString(),
		Name:        uuid.NewString(),
		Value:       "s:" + service + ";" + "ns:" + namespace,
		NamespaceId: namespace,
		ServiceId:   service,
		Active:      active,
	}
}

func CreatePropActive(service, namespace string) utils.DcPropertyDto {
	return CreateProp(service, namespace, true)
}

func CreateNs(name string) NamespaceDto {
	return NamespaceDto{
		Id:   uuid.NewString(),
		Name: name,
	}
}

func CreateService(name string) ServiceDto {
	return ServiceDto{
		Id:   uuid.NewString(),
		Name: name,
	}
}

func MergeSlices[T any](args ...[]T) []T {
	result := make([]T, 0)
	for _, s := range args {
		result = append(result, s...)
	}
	return result
}

func (d Datasets) OnlyAllServicesDefault() Dataset {
	allProps := []utils.DcPropertyDto{
		CreatePropActive("", ""),
		CreatePropActive("", ""),
		CreatePropActive("", ""),
	}
	return Dataset{
		ImportExportAnswer{
			Namespaces: []NamespaceDto{},
			Services:   []ServiceDto{},
			Props:      allProps,
		},
		[]Assertion{
			{"", "", allProps},
			{"some", "", allProps},
			{"some2", "", allProps},
			{"", "some", emptyAns},
			{"some", "some", emptyAns},
			{"some", "some2", emptyAns},
		},
		[]string{},
	}
}

func (d Datasets) OnlyAllServicesWithNamespace() Dataset {
	ns := CreateNs("ns")
	allProps := []utils.DcPropertyDto{
		CreatePropActive("", ns.Id),
		CreatePropActive("", ns.Id),
		CreatePropActive("", ns.Id),
	}
	return Dataset{
		ImportExportAnswer{
			Namespaces: []NamespaceDto{ns},
			Services:   []ServiceDto{},
			Props:      allProps,
		},
		[]Assertion{
			{"", "", emptyAns},
			{"", "ns", allProps},
			{"", "ns2", emptyAns},
			{"s", "", emptyAns},
			{"s", "ns", allProps},
			{"s", "ns2", emptyAns},
			{"s2", "", emptyAns},
			{"s2", "ns", allProps},
			{"s2", "ns2", emptyAns},
		},
		[]string{},
	}
}

func (d Datasets) OneServiceNoNamespace() Dataset {
	s := CreateService("s")
	allProps := []utils.DcPropertyDto{
		CreatePropActive(s.Id, ""),
		CreatePropActive(s.Id, ""),
		CreatePropActive(s.Id, ""),
	}
	return Dataset{
		ImportExportAnswer{
			Namespaces: []NamespaceDto{},
			Services:   []ServiceDto{s},
			Props:      allProps,
		},
		[]Assertion{
			{"", "", emptyAns},
			{"", "ns", emptyAns},
			{"", "ns2", emptyAns},
			{"s", "", allProps},
			{"s", "ns", emptyAns},
			{"s", "ns2", emptyAns},
			{"s2", "", emptyAns},
			{"s2", "ns", emptyAns},
			{"s2", "ns2", emptyAns},
		},
		[]string{},
	}
}

func (d Datasets) ManyServicesOneNs() Dataset {
	s := CreateService("s")
	s2 := CreateService("s2")
	ns := CreateNs("ns")
	s1Props := []utils.DcPropertyDto{
		CreatePropActive(s.Id, ns.Id),
		CreatePropActive(s.Id, ns.Id),
	}
	s2Props := []utils.DcPropertyDto{
		CreatePropActive(s2.Id, ns.Id),
		CreatePropActive(s2.Id, ns.Id),
	}
	allServicesProps := []utils.DcPropertyDto{
		CreatePropActive("", ns.Id),
		CreatePropActive("", ns.Id),
	}
	allProps := MergeSlices(s1Props, s2Props, allServicesProps)
	return Dataset{
		ImportExportAnswer{
			Namespaces: []NamespaceDto{ns},
			Services:   []ServiceDto{s, s2},
			Props:      allProps,
		},
		[]Assertion{
			{"", "", emptyAns},
			{"s", "", emptyAns},
			{"s2", "", emptyAns},
			{"s3", "", emptyAns},
			{"", "ns2", emptyAns},
			{"s", "ns2", emptyAns},
			{"s2", "ns2", emptyAns},
			{"s3", "ns2", emptyAns},
			{"", "ns", allServicesProps},
			{"s", "ns", MergeSlices(s1Props, allServicesProps)},
			{"s2", "ns", MergeSlices(s2Props, allServicesProps)},
			{"s3", "ns", allServicesProps},
		},
		[]string{},
	}
}

func (d Datasets) ManyServicesManyNs() Dataset {
	s := CreateService("s")
	s2 := CreateService("s2")
	ns := CreateNs("ns")
	ns2 := CreateNs("ns2")
	s1n1Props := []utils.DcPropertyDto{
		CreatePropActive(s.Id, ns.Id),
		CreatePropActive(s.Id, ns.Id),
	}
	s1n2Props := []utils.DcPropertyDto{
		CreatePropActive(s.Id, ns2.Id),
		CreatePropActive(s.Id, ns2.Id),
	}
	s2n1Props := []utils.DcPropertyDto{
		CreatePropActive(s2.Id, ns.Id),
		CreatePropActive(s2.Id, ns.Id),
	}
	s2n2Props := []utils.DcPropertyDto{
		CreatePropActive(s2.Id, ns2.Id),
		CreatePropActive(s2.Id, ns2.Id),
	}
	s1Props := []utils.DcPropertyDto{
		CreatePropActive(s.Id, ""),
		CreatePropActive(s.Id, ""),
	}
	s2Props := []utils.DcPropertyDto{
		CreatePropActive(s2.Id, ""),
		CreatePropActive(s2.Id, ""),
	}
	ns1Props := []utils.DcPropertyDto{
		CreatePropActive("", ns.Id),
		CreatePropActive("", ns.Id),
	}
	ns2Props := []utils.DcPropertyDto{
		CreatePropActive("", ns2.Id),
		CreatePropActive("", ns2.Id),
	}
	props := []utils.DcPropertyDto{
		CreatePropActive("", ""),
		CreatePropActive("", ""),
	}
	allProps := MergeSlices(s1n1Props, s1n2Props, s2n1Props, s2n2Props, s1Props, s2Props, ns1Props, ns2Props, props)
	return Dataset{
		ImportExportAnswer{
			Namespaces: []NamespaceDto{ns, ns2},
			Services:   []ServiceDto{s, s2},
			Props:      allProps,
		},
		[]Assertion{
			{"", "", MergeSlices(props)},
			{"s", "", MergeSlices(props, s1Props)},
			{"s2", "", MergeSlices(props, s2Props)},
			{"s3", "", props},
			{"", "ns", ns1Props},
			{"s", "ns", MergeSlices(ns1Props, s1n1Props)},
			{"s2", "ns", MergeSlices(ns1Props, s2n1Props)},
			{"s3", "ns", ns1Props},
			{"", "ns2", ns2Props},
			{"s", "ns2", MergeSlices(ns2Props, s1n2Props)},
			{"s2", "ns2", MergeSlices(ns2Props, s2n2Props)},
			{"s3", "ns2", ns2Props},
			{"", "ns3", emptyAns},
			{"s", "ns3", emptyAns},
			{"s2", "ns3", emptyAns},
			{"s3", "ns3", emptyAns},
		},
		[]string{},
	}
}

func (d Datasets) SamePropWithAndWithoutNs_withService() Dataset {
	s := CreateService("s")
	ns := CreateNs("ns")
	propWithNs := CreatePropFull("id1", "K", "V", s.Id, ns.Id, true)
	propWithoutNs := CreatePropFull("id2", "K", "V", s.Id, "", true)
	return Dataset{
		ImportExportAnswer{
			Namespaces: []NamespaceDto{ns},
			Services:   []ServiceDto{s},
			Props:      []utils.DcPropertyDto{propWithNs, propWithoutNs},
		},
		[]Assertion{
			{"s", "", []utils.DcPropertyDto{propWithoutNs}},
			{"s", "ns", []utils.DcPropertyDto{propWithNs}},
		},
		[]string{
			"if ns not specified - prop without ns overrides prop with ns; with service",
			"if ns specified - prop with ns overrides prop without ns; with service",
		},
	}
}

func (d Datasets) SamePropWithAndWithoutNs_withoutService() Dataset {
	ns := CreateNs("ns")
	propWithNs := CreatePropFull("id1", "K", "V", "", ns.Id, true)
	propWithoutNs := CreatePropFull("id2", "K", "V", "", "", true)
	return Dataset{
		ImportExportAnswer{
			Namespaces: []NamespaceDto{ns},
			Services:   []ServiceDto{},
			Props:      []utils.DcPropertyDto{propWithNs, propWithoutNs},
		},
		[]Assertion{
			{"", "", []utils.DcPropertyDto{propWithoutNs}},
			{"", "ns", []utils.DcPropertyDto{propWithNs}},
		},
		[]string{
			"if ns not specified - prop without ns overrides prop with ns; without service",
			"if ns specified - prop with ns overrides prop without ns; without service",
		},
	}
}

func (d Datasets) SamePropWithAndWithoutService_withNs() Dataset {
	s := CreateService("s")
	ns := CreateNs("ns")
	propWithService := CreatePropFull("id1", "K", "V", s.Id, ns.Id, true)
	propWithoutService := CreatePropFull("id2", "K", "V", "", ns.Id, true)
	return Dataset{
		ImportExportAnswer{
			Namespaces: []NamespaceDto{ns},
			Services:   []ServiceDto{s},
			Props:      []utils.DcPropertyDto{propWithService, propWithoutService},
		},
		[]Assertion{
			{"", "ns", []utils.DcPropertyDto{propWithoutService}},
			{"s", "ns", []utils.DcPropertyDto{propWithService}},
		},
		[]string{
			"if service not specified - prop with all services (empty service) overrides prop with service; with ns",
			"if service specified - prop with service overrides prop with all services (empty service); with ns",
		},
	}
}

func (d Datasets) SamePropWithAndWithoutService_withoutNs() Dataset {
	s := CreateService("s")
	propWithService := CreatePropFull("id1", "K", "V", s.Id, "", true)
	propWithoutService := CreatePropFull("id2", "K", "V", "", "", true)
	return Dataset{
		ImportExportAnswer{
			Namespaces: []NamespaceDto{},
			Services:   []ServiceDto{s},
			Props:      []utils.DcPropertyDto{propWithService, propWithoutService},
		},
		[]Assertion{
			{"", "", []utils.DcPropertyDto{propWithoutService}},
			{"s", "", []utils.DcPropertyDto{propWithService}},
		},
		[]string{
			"if service not specified - prop with all services overrides prop with service; without ns",
			"if service specified - prop with service overrides prop with all services (empty service); without ns",
		},
	}
}
