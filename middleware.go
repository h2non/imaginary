package main

import (
	"fmt"
	"github.com/daaku/go.httpgzip"
	"github.com/rs/cors"
	"github.com/throttled/throttled"
	"gopkg.in/h2non/bimg.v0"
	"net/http"
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
		return Middleware(imageControllerDispatcher(o, Operation(fn)), o)
	}
}

func throttle(next http.Handler, o ServerOptions) http.Handler {
	th := throttled.Interval(throttled.PerSec(o.Concurrency), o.Burst, &throttled.VaryBy{Method: true}, o.Burst)
	return th.Throttle(next)
}

func validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" && r.Method != "POST" {
			ErrorReply(w, ErrMethodNotAllowed)
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
			ErrorReply(w, ErrInvalidApiKey)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func defaultHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", fmt.Sprintf("imaginary %s (bimg %s)", Version, bimg.Version))
		next.ServeHTTP(w, r)
	})
}
