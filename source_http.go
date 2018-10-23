package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
	url, err := parseURL(req)
	if err != nil {
		return nil, ErrInvalidImageURL
	}
	if shouldRestrictOrigin(url, s.Config.AllowedOrigins) {
		return nil, fmt.Errorf("Not allowed remote URL origin: %s", url.Host)
	}
	return s.fetchImage(url, req)
}

func (s *HttpImageSource) fetchImage(url *url.URL, ireq *http.Request) ([]byte, error) {
	// Check remote image size by fetching HTTP Headers
	if s.Config.MaxAllowedSize > 0 {
		req := newHTTPRequest(s, ireq, "HEAD", url)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("Error fetching image http headers: %v", err)
		}
		res.Body.Close()
		if res.StatusCode < 200 && res.StatusCode > 206 {
			return nil, fmt.Errorf("Error fetching image http headers: (status=%d) (url=%s)", res.StatusCode, req.URL.String())
		}

		contentLength, _ := strconv.Atoi(res.Header.Get("Content-Length"))
		if contentLength > s.Config.MaxAllowedSize {
			return nil, fmt.Errorf("Content-Length %d exceeds maximum allowed %d bytes", contentLength, s.Config.MaxAllowedSize)
		}
	}

	// Perform the request using the default client
	req := newHTTPRequest(s, ireq, "GET", url)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error downloading image: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Error downloading image: (status=%d) (url=%s)", res.StatusCode, req.URL.String())
	}

	// Read the body
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Unable to create image from response body: %s (url=%s)", req.URL.String(), err)
	}
	return buf, nil
}

func (s *HttpImageSource) setAuthorizationHeader(req *http.Request, ireq *http.Request) {
	auth := s.Config.Authorization
	if auth == "" {
		auth = ireq.Header.Get("X-Forward-Authorization")
	}
	if auth == "" {
		auth = ireq.Header.Get("Authorization")
	}
	if auth != "" {
		req.Header.Set("Authorization", auth)
	}
}

func parseURL(request *http.Request) (*url.URL, error) {
	queryUrl := request.URL.Query().Get("url")
	return url.Parse(queryUrl)
}

func newHTTPRequest(s *HttpImageSource, ireq *http.Request, method string, url *url.URL) *http.Request {
	req, _ := http.NewRequest(method, url.String(), nil)
	req.Header.Set("User-Agent", "imaginary/"+Version)
	req.URL = url

	// Forward auth header to the target server, if necessary
	if s.Config.AuthForwarding || s.Config.Authorization != "" {
		s.setAuthorizationHeader(req, ireq)
	}

	return req
}

func shouldRestrictOrigin(url *url.URL, origins []*url.URL) bool {
	if len(origins) == 0 {
		return false
	}

	for _, origin := range origins {
		if origin.Host == url.Host {
			return false
		}

		if origin.Host[0:2] == "*." {

			// Testing if "*.example.org" matches "example.org"
			if url.Host == origin.Host[2:] {
				return false
			}

			// Testing if "*.example.org" matches "foo.example.org"
			if strings.HasSuffix(url.Host, origin.Host[1:]) {
				return false
			}
		}
	}

	return true
}

func init() {
	RegisterSource(ImageSourceTypeHttp, NewHttpImageSource)
}
