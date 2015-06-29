package main

import (
	"fmt"
	"github.com/PuerkitoBio/throttled"
	"github.com/daaku/go.httpgzip"
	"github.com/rs/cors"
	"gopkg.in/h2non/bimg.v0"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
)

func Middleware(fn func(http.ResponseWriter, *http.Request), o ServerOptions) http.Handler {
	next := http.Handler(http.HandlerFunc(fn))

	if o.Concurrency > 0 {
		next = throttle(next, o)
	}
	if o.Gzip {
		next = httpgzip.NewHandler(next)
	}
	if o.CORS {
		next = cors.Default().Handler(next)
	}
	if o.ApiKey != "" {
		next = validateApiKey(next, o.ApiKey)
	}

	return validate(defaultHeaders(next))
}

func ImageMiddleware(o ServerOptions) func(Operation) http.Handler {
	return func(fn Operation) http.Handler {
		return Middleware(func(w http.ResponseWriter, r *http.Request) {
			validMethod := r.Method != "POST" || (r.Method == "GET" && o.Mount != "")

			if validMethod == false {
				ErrorReply(w, "Method not allowed for this endpoint", NOT_ALLOWED)
				return
			}

			file := r.URL.Query().Get("file")
			if file != "" && r.Method == "GET" {
				file = path.Clean(path.Join(o.Mount, file))
				if strings.HasPrefix(file, o.Mount) == false {
					ErrorReply(w, "Invalid local file path", BAD_REQUEST)
					return
				}

				buf, err := ioutil.ReadFile(file)
				if err != nil {
					ErrorReply(w, "Invalid local file path", BAD_REQUEST)
					return
				}

				imageController(w, r, buf, Operation(fn))
				return
			}

			buf, err := readBody(r)
			if err != nil {
				ErrorReply(w, "Cannot read the payload: "+err.Error(), BAD_REQUEST)
				return
			}

			imageController(w, r, buf, Operation(fn))
		}, o)
	}
}

func throttle(next http.Handler, o ServerOptions) http.Handler {
	th := throttled.Interval(throttled.PerSec(o.Concurrency), o.Burst, &throttled.VaryBy{Method: true}, o.Burst)
	return th.Throttle(next)
}

func validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" && r.Method != "POST" {
			ErrorReply(w, "Method not allowed: "+r.Method, NOT_ALLOWED)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func validateApiKey(next http.Handler, validKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("API-Key")
		if key == "" {
			key = r.URL.Query().Get("key")
		}

		if key != validKey {
			ErrorReply(w, "Invalid or missing API key", UNAUTHORIZED)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func defaultHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", fmt.Sprintf("imaginary %s (using bimg %s)", Version, bimg.Version))
		next.ServeHTTP(w, r)
	})
}
