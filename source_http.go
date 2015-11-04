package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const ImageSourceTypeHttp ImageSourceType = "http"

type HttpImageSource struct {
	Config *SourceConfig
}

func NewHttpImageSource(config *SourceConfig) ImageSource {
	return &HttpImageSource{config}
}

func (s *HttpImageSource) Matches(r *http.Request) bool {
	return r.Method == "GET" && r.URL.Query().Get("url") != ""
}

func (s *HttpImageSource) GetImage(req *http.Request) ([]byte, error) {
	url, err := s.parseURL(req)
	if err != nil {
		return nil, ErrInvalidImageURL
	}
	return s.fetchImage(url)
}

func (s *HttpImageSource) fetchImage(url *url.URL) ([]byte, error) {
	req := s.newHttpRequest(url)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error downloading image: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Error downloading image: (status=%d) (url=%s)", res.StatusCode, req.URL.RequestURI())
	}

	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Unable to create image from response body: %s (url=%s)", req.URL.RequestURI(), err)
	}
	return buf, nil
}

func (s *HttpImageSource) parseURL(request *http.Request) (*url.URL, error) {
	queryUrl := request.URL.Query().Get("url")
	return url.Parse(queryUrl)
}

func (s *HttpImageSource) newHttpRequest(url *url.URL) *http.Request {
	req, _ := http.NewRequest("GET", url.RequestURI(), nil)
	req.Header.Set("User-Agent", "imaginary")
	req.URL = url
	return req
}

func init() {
	RegisterSource(ImageSourceTypeHttp, NewHttpImageSource)
}
