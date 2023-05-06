package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type (
	Client struct {
		Http        http.Client
		BaseUrl     string
		Token       string
		LogRequest  bool
		LogResponse bool
	}
)

var NotFound = errors.New("resource not found")
var Conflict = errors.New("resource conflict")

func JSON() map[string]string {
	return map[string]string{"Content-Type": "application/json"}
}

func (k *Client) DoRaw(method string, path string, headers map[string]string, requestBody []byte) ([]byte, error) {
	fullUrl := k.BaseUrl + path
	request, err := http.NewRequest(method, fullUrl, nil)
	if err != nil {
		return nil, err
	}

	if k.Token != "" {
		request.Header.Add("Authorization", "Bearer "+k.Token)
	}
	if headers != nil {
		for k, v := range headers {
			request.Header.Add(k, v)
		}
	}

	request.Body = io.NopCloser(bytes.NewReader(requestBody))
	response, err := k.Http.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)

	switch response.StatusCode {
	case 404:
		return responseBody, NotFound
	case 409:
		return responseBody, Conflict
	default:
		if response.StatusCode >= 400 {
			return responseBody, fmt.Errorf("got status code %d from REST API (%s %s)", response.StatusCode, method, fullUrl)
		}
		return responseBody, err
	}
}

func (k *Client) Do(method string, path string, headers map[string]string, request interface{}, result interface{}) error {
	var requestBody []byte
	var err error
	var data []byte

	if request != nil {
		requestBody, err = json.Marshal(request)
		if k.LogRequest {
			log.Printf("formatted %s request to %s as %s\n", method, k.BaseUrl+path, string(requestBody))
		}
	}
	if err != nil {
		return err
	}
	data, err = k.DoRaw(method, path, headers, requestBody)
	if err != nil {
		return err
	}
	if result != nil {
		if k.LogResponse {
			log.Printf("received response from %s: %s\n", k.BaseUrl+path, string(data))
		}
		return json.Unmarshal(data, result)
	}
	return nil
}
