package screws

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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

//MakeTimestamp ...
func (at *APITester) MakeTimestamp(datetime ...string) []int64 {
	var timestamps []int64
	for _, v := range datetime {
		time, err := time.ParseInLocation("2006-01-02 15:04:05", v, time.Local)
		if err != nil {
			log.Fatal(err)
		}
		timestamps = append(timestamps, time.Unix())
	}
	return timestamps
}
