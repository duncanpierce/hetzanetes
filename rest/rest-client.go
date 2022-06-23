package rest

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"k8s.io/apimachinery/pkg/util/json"
	"log"
	"net/http"
)

type (
	Client struct {
		Http    http.Client
		BaseUrl string
		Token   string
	}
)

var NotFound = errors.New("resource not found")

func JSON() map[string]string {
	return map[string]string{"Content-Type": "application/json"}
}

func (k *Client) DoRaw(method string, path string, headers map[string]string, requestBody []byte) ([]byte, error) {
	request, err := http.NewRequest(method, k.BaseUrl+path, nil)
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
	if response.StatusCode == 404 {
		return responseBody, NotFound
	}
	if response.StatusCode >= 400 {
		return responseBody, fmt.Errorf("got status code %d from Kubernetes API", response.StatusCode)
	}
	return responseBody, err
}

func (k *Client) Do(method string, path string, headers map[string]string, request interface{}, result interface{}) error {
	var requestBody []byte
	var err error
	var data []byte

	if request != nil {
		requestBody, err = json.Marshal(request)
		log.Printf("formatted %s to %s as %s\n", method, k.BaseUrl+path, string(requestBody))
	}
	if err != nil {
		return err
	}
	data, err = k.DoRaw(method, path, headers, requestBody)
	if err != nil {
		return err
	}
	if result != nil {
		return json.Unmarshal(data, result)
	}
	return nil
}
