package screws

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//APITester ...
type APITester struct {
	URL           string
	Method        string
	RequestHeader map[string]string
	RequestData   interface{}
	ResponseData  interface{}
}

//LoadHTTPRequestWithJSON ...
func (at *APITester) LoadHTTPRequestWithJSON() error {
	requestBody, err := json.Marshal(at.RequestData)
	if err != nil {
		return err
	}
	if at.ResponseData, err = loadHTTPRequest(at.Method, at.URL, requestBody, at.RequestHeader); err != nil {
		return err
	}
	return nil
}

//RunTest ...
func (at *APITester) RunTest(method string, url string, requestHeader map[string]string, data interface{}) {
	at.Method = method
	at.URL = url
	at.RequestHeader = requestHeader
	at.RequestData = data
	if err := at.LoadHTTPRequestWithJSON(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Request:%v %v %v\nResponse: %v\n", at.Method, at.URL, at.RequestHeader, at.ResponseData)
}

//loadHTTPRequest ...
func loadHTTPRequest(method string, url string, body []byte, head map[string]string) (string, error) {
	request, err := http.NewRequest(method, url, bytes.NewReader(body))
	for k, v := range head {
		request.Header.Set(k, v)
	}
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	responseBodyByte, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v %v %v", response.Status, response.Header, string(responseBodyByte)), nil
}
