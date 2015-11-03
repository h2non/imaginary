package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	ImageSourceTypeHttp ImageSourceType = "http"
)

type HttpImageSource struct {
	Config *SourceConfig
}

func NewHttpImageSourceWithConfig(config *SourceConfig) ImageSource {
	return &HttpImageSource{config}
}

func GetImageFromReader(buffer io.Reader) ([]byte, error) {
	return ioutil.ReadAll(buffer)
}

func (s *HttpImageSource) GetImage(request *http.Request) ([]byte, error) {
	URL, err := s.parseURL(request)
	if err != nil {
		return nil, err
	}

	httpRequest := s.newHttpRequest(URL)
	httpResponse, err := http.DefaultClient.Do(httpRequest)
	defer httpResponse.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Error downloading image: %v", err)
	}
	if httpResponse.StatusCode != 200 {
		return nil, fmt.Errorf("Error downloading image (url=%s)", httpRequest.URL.RequestURI())
	}

	body, err := GetImageFromReader(httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("Unable to create image from response body: %v (url=%s)", string(body), httpRequest.URL.RequestURI())
	}

	return body, nil
}

func (s *HttpImageSource) parseURL(request *http.Request) (*url.URL, error) {
	queryUrl := request.URL.Query().Get("url")
	return url.Parse(queryUrl)
}

func (s *HttpImageSource) newHttpRequest(url *url.URL) *http.Request {
	httpRequest, _ := http.NewRequest("GET", url.RequestURI(), nil)
	httpRequest.URL = url
	return httpRequest
}

func init() {
	RegisterSource(ImageSourceTypeHttp, NewHttpImageSourceWithConfig)
}
