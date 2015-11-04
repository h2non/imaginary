package main

import (
	"fmt"
	"github.com/daaku/go.httpgzip"
	"github.com/rs/cors"
	"gopkg.in/h2non/bimg.v0"
	"gopkg.in/throttled/throttled.v2"
	"gopkg.in/throttled/throttled.v2/store/memstore"
	"net/http"
	"time"
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
		next = authorizeClient(next, o.ApiKey)
	}
	if o.HttpCacheTtl >= 0 {
		next = defineCacheHeaders(next, o.HttpCacheTtl)
	}

	return validate(defaultHeaders(next))
}

func ImageMiddleware(o ServerOptions) func(Operation) http.Handler {
	return func(fn Operation) http.Handler {
		return validateImage(Middleware(imageController(o, Operation(fn)), o), o)
	}
}

func throttleError(err error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "throttle error: "+err.Error(), http.StatusInternalServerError)
	})
}

func throttle(next http.Handler, o ServerOptions) http.Handler {
	store, err := memstore.New(65536)
	if err != nil {
		return throttleError(err)
	}

	quota := throttled.RateQuota{throttled.PerSec(o.Concurrency), o.Burst}
	rateLimiter, err := throttled.NewGCRARateLimiter(store, quota)
	if err != nil {
		return throttleError(err)
	}

	httpRateLimiter := throttled.HTTPRateLimiter{
		RateLimiter: rateLimiter,
		VaryBy:      &throttled.VaryBy{Method: true},
	}

	return httpRateLimiter.RateLimit(next)
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

func validateImage(next http.Handler, o ServerOptions) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if r.Method == "GET" && (path == "/" || path == "/health" || path == "/form") {
			next.ServeHTTP(w, r)
			return
		}

		if r.Method == "GET" && o.Mount == "" && o.EnableURLSource == false {
			ErrorReply(w, ErrMethodNotAllowed)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func authorizeClient(next http.Handler, validKey string) http.Handler {
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

func defineCacheHeaders(next http.Handler, ttl int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer next.ServeHTTP(w, r)
		if r.Method != "GET" {
			return
		}

		var cacheControl string
		if ttl == 0 {
			cacheControl = "private, no-cache, no-store, must-revalidate"
		} else {
			cacheControl = fmt.Sprintf("public, s-maxage: %d, max-age: %d, no-transform", ttl, ttl)
		}

		ttlDiff := time.Duration(ttl) * time.Second
		expires := time.Now().Add(ttlDiff)
		w.Header().Add("Expires", expires.Format(time.RFC1123))
		w.Header().Add("Cache-Control", cacheControl)
	})
}
