package util

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"sync"
)

var httpClient *http.Client
var once sync.Once

func GetWeb(url string) (io.ReadCloser, error) {
	response, err := getHttpClient().Get(url)
	if err != nil {
		fmt.Println("error downloading web page")

		return nil, err
	}

	if response.StatusCode < 200 && response.StatusCode > 299 {
		fmt.Println("error downloading web page, status is not 200")

		return nil, err
	}

	return response.Body, nil
}

func getHttpClient() *http.Client {
	once.Do(func() {
		config := &tls.Config{InsecureSkipVerify: true}
		transport := &http.Transport{TLSClientConfig: config}
		httpClient = &http.Client{Transport: transport}
	})

	return httpClient
}
