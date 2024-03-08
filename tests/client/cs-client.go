package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	utils "github.com/logotipiwe/dc_go_utils/src"
	"net/http"
	"net/url"
	"os"
	"tests/datasets"
)

type CsClient struct {
}

func NewCsClient() *CsClient {
	return &CsClient{}
}

func (c CsClient) getCsUrl() string {
	return os.Getenv("CS_URL")
}

func (c CsClient) getMToken() string {
	return os.Getenv("M_TOKEN")
}

func (c CsClient) GetProps(service, namespace string) ([]utils.DcPropertyDto, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", c.getCsUrl()+"/api/get-config", nil)
	query := req.URL.Query()
	query.Add("mToken", c.getMToken())
	query.Add("service", service)
	query.Add("namespace", namespace)
	req.URL.RawQuery = query.Encode()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	answer := make([]utils.DcPropertyDto, 0)
	err = json.NewDecoder(res.Body).Decode(&answer)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return answer, nil
}

func (c CsClient) ImportConfig(data datasets.ImportExportAnswer) error {
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(data)
	request, err := http.NewRequest("POST", c.getCsUrl()+"/api/import", buffer)
	params := url.Values{}
	params.Add("mToken", c.getMToken())
	//params.Add("service", os.Getenv("SERVICE_NAME"))
	//params.Add("namespace", os.Getenv("NAMESPACE"))

	request.URL.RawQuery = params.Encode()
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return errors.New("Got status " + res.Status)
	}
	//defer res.Body.Close()
	//var answer []utils.DcPropertyDto
	//err = json.NewDecoder(res.Body).Decode(&answer)
	return nil
}

func (c CsClient) ExportConfig() (*datasets.ImportExportAnswer, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", c.getCsUrl()+"/api/export", nil)
	query := req.URL.Query()
	query.Add("mToken", c.getMToken())
	req.URL.RawQuery = query.Encode()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	answer := datasets.ImportExportAnswer{}
	err = json.NewDecoder(res.Body).Decode(&answer)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &answer, nil
}
