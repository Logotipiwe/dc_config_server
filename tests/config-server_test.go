package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sort"
	"strconv"
	"testing"
	csClient "tests/client"
	"tests/datasets"
)

func importAndCheck(t *testing.T, client *csClient.CsClient, dataset datasets.Dataset) {
	err := client.ImportConfig(dataset.Data)
	assert.Nil(t, err)
	if err != nil {
		return
	}
	for i, assertion := range dataset.Assertions {
		fmt.Println("Checking entry " + strconv.Itoa(i+1))
		props, err := client.GetProps(assertion.Service, assertion.Namespace)

		testFunc := func(t *testing.T) {
			assert.Nil(t, err)
			sort.Slice(props, func(i, j int) bool {
				return props[i].Id < props[j].Id
			})
			sort.Slice(assertion.Answer, func(i, j int) bool {
				return assertion.Answer[i].Id < assertion.Answer[j].Id
			})
			assert.Equal(t, assertion.Answer, props)
		}

		if len(dataset.TestNames) > 0 {
			t.Run(dataset.TestNames[i], testFunc)
		} else {
			t.Run(strconv.Itoa(i)+" test case", testFunc)
		}
	}
}

func TestConfigServerApi(t *testing.T) {
	client := csClient.NewCsClient()
	ds := datasets.Datasets{}
	t.Run("all services no ns", func(t *testing.T) {
		importAndCheck(t, client, ds.OnlyAllServicesDefault())
	})
	t.Run("all services some ns", func(t *testing.T) {
		importAndCheck(t, client, ds.OnlyAllServicesWithNamespace())
	})
	t.Run("one service no ns", func(t *testing.T) {
		importAndCheck(t, client, ds.OneServiceNoNamespace())
	})
	t.Run("many services one ns", func(t *testing.T) {
		importAndCheck(t, client, ds.ManyServicesOneNs())
	})
	t.Run("many services many ns", func(t *testing.T) {
		importAndCheck(t, client, ds.ManyServicesManyNs())
	})

	t.Run("ns overriding; with service", func(t *testing.T) {
		importAndCheck(t, client, ds.SamePropWithAndWithoutNs_withService())
	})
	t.Run("ns overriding; without service", func(t *testing.T) {
		importAndCheck(t, client, ds.SamePropWithAndWithoutNs_withoutService())
	})

	t.Run("service overriding; with ns", func(t *testing.T) {
		importAndCheck(t, client, ds.SamePropWithAndWithoutService_withNs())
	})
	t.Run("service overriding; without ns", func(t *testing.T) {
		importAndCheck(t, client, ds.SamePropWithAndWithoutService_withoutNs())
	})
}
