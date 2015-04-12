package main

import (
	"github.com/daaku/go.httpgzip"
	"github.com/rs/cors"
	"net/http"
)

func Middleware(fn func(http.ResponseWriter, *http.Request), o ServerOptions) http.Handler {
	next := http.Handler(http.HandlerFunc(fn))

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
			imageController(w, r, Operation(fn))
		}, o)
	}
}

func validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" && r.Method != "POST" {
			errorResponse(w, "Method not allowed: "+r.Method, NOT_ALLOWED)
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
			errorResponse(w, "Invalid or missing API key", UNAUTHORIZED)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func defaultHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "imaginary "+Version)
		next.ServeHTTP(w, r)
	})
}
