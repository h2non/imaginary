package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const ImageSourceTypeHTTP ImageSourceType = "http"
const URLQueryKey = "url"

type HTTPImageSource struct {
	Config *SourceConfig
}

func NewHTTPImageSource(config *SourceConfig) ImageSource {
	return &HTTPImageSource{config}
}

func (s *HTTPImageSource) Matches(r *http.Request) bool {
	return r.Method == http.MethodGet && r.URL.Query().Get(URLQueryKey) != ""
}

func (s *HTTPImageSource) GetImage(req *http.Request) ([]byte, error) {
	u, err := parseURL(req)
	if err != nil {
		return nil, ErrInvalidImageURL
	}
	if shouldRestrictOrigin(u, s.Config.AllowedOrigins) {
		return nil, fmt.Errorf("not allowed remote URL origin: %s%s", u.Host, u.Path)
	}
	return s.fetchImage(u, req)
}

func (s *HTTPImageSource) fetchImage(murl *url.URL, ireq *http.Request) ([]byte, error) {

	queryURL := murl.String()

	if strings.Contains(queryURL, "%") {
		var err error
		queryURL, err = url.QueryUnescape(queryURL)
		if err != nil {
			fmt.Printf("failed to unesacpe url: %v", err)
		}
		fmt.Printf("queryURL unescape: %s\n", queryURL)

		murl, _ = url.Parse(queryURL)
	}

	// Check remote image size by fetching HTTP Headers
	if s.Config.MaxAllowedSize > 0 {
		req := newHTTPRequest(s, ireq, http.MethodHead, murl)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error fetching remote http image headers: %v", err)
		}
		_ = res.Body.Close()
		if res.StatusCode < 200 && res.StatusCode > 206 {
			return nil, NewError(fmt.Sprintf("error fetching remote http image headers: (status=%d) (url=%s)", res.StatusCode, req.URL.String()), res.StatusCode)
		}

		contentLength, _ := strconv.Atoi(res.Header.Get("Content-Length"))
		if contentLength > s.Config.MaxAllowedSize {
			return nil, fmt.Errorf("Content-Length %d exceeds maximum allowed %d bytes", contentLength, s.Config.MaxAllowedSize)
		}
	}

	// Perform the request using the default client
	req := newHTTPRequest(s, ireq, http.MethodGet, murl)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching remote http image: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, NewError(fmt.Sprintf("error fetching remote http image: (status=%d) (url=%s)", res.StatusCode, req.URL.String()), res.StatusCode)
	}

	// Read the body
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to create image from response body: %s (url=%s)", req.URL.String(), err)
	}
	return buf, nil
}

func (s *HTTPImageSource) setAuthorizationHeader(req *http.Request, ireq *http.Request) {
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

func (s *HTTPImageSource) setForwardHeaders(req *http.Request, ireq *http.Request) {
	headers := s.Config.ForwardHeaders
	for _, header := range headers {
		if _, ok := ireq.Header[header]; ok {
			req.Header.Set(header, ireq.Header.Get(header))
		}
	}
}

func parseURL(request *http.Request) (*url.URL, error) {
	queryURL := request.URL.Query().Get("url")

	if strings.Contains(queryURL, "%") {
		var err error
		queryURL, err = url.QueryUnescape(queryURL)
		if err != nil {
			fmt.Printf("failed to unescape url: %v", err)
		}
		//fmt.Printf("queryURL unescape: %s\n", queryURL)
	}

	return url.Parse(queryURL)
}

func newHTTPRequest(s *HTTPImageSource, ireq *http.Request, method string, url *url.URL) *http.Request {
	req, _ := http.NewRequest(method, url.String(), nil)
	req.Header.Set("User-Agent", "imaginary/"+Version)
	req.URL = url

	if len(s.Config.ForwardHeaders) != 0 {
		s.setForwardHeaders(req, ireq)
	}

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
			if strings.HasPrefix(url.Path, origin.Path) {
				return false
			}
		}

		if origin.Host[0:2] == "*." {
			// Testing if "*.example.org" matches "example.org"
			if url.Host == origin.Host[2:] {
				if strings.HasPrefix(url.Path, origin.Path) {
					return false
				}
			}

			// Testing if "*.example.org" matches "foo.example.org"
			if strings.HasSuffix(url.Host, origin.Host[1:]) {
				if strings.HasPrefix(url.Path, origin.Path) {
					return false
				}
			}
		}
	}

	return true
}

func init() {
	RegisterSource(ImageSourceTypeHTTP, NewHTTPImageSource)
}
