package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/h2non/bimg"
	"github.com/rs/cors"
	"gopkg.in/throttled/throttled.v2"
	"gopkg.in/throttled/throttled.v2/store/memstore"
)

func Middleware(fn func(http.ResponseWriter, *http.Request), o ServerOptions) http.Handler {
	next := http.Handler(http.HandlerFunc(fn))

	if len(o.Endpoints) > 0 {
		next = filterEndpoint(next, o)
	}
	if o.Concurrency > 0 {
		next = throttle(next, o)
	}
	if o.CORS {
		next = cors.Default().Handler(next)
	}
	if o.APIKey != "" {
		next = authorizeClient(next, o)
	}
	if o.HTTPCacheTTL >= 0 {
		next = setCacheHeaders(next, o.HTTPCacheTTL)
	}

	return validate(defaultHeaders(next), o)
}

func ImageMiddleware(o ServerOptions) func(Operation) http.Handler {
	return func(fn Operation) http.Handler {
		handler := validateImage(Middleware(imageController(o, fn), o), o)

		if o.EnableURLSignature {
			return validateURLSignature(handler, o)
		}

		return handler
	}
}

func filterEndpoint(next http.Handler, o ServerOptions) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if o.Endpoints.IsValid(r) {
			next.ServeHTTP(w, r)
			return
		}
		ErrorReply(r, w, ErrNotImplemented, o)
	})
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

	quota := throttled.RateQuota{MaxRate: throttled.PerSec(o.Concurrency), MaxBurst: o.Burst}
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

func validate(next http.Handler, o ServerOptions) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			ErrorReply(r, w, ErrMethodNotAllowed, o)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func validateImage(next http.Handler, o ServerOptions) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if r.Method == http.MethodGet && isPublicPath(path) {
			next.ServeHTTP(w, r)
			return
		}

		if r.Method == http.MethodGet && o.Mount == "" && !o.EnableURLSource {
			ErrorReply(r, w, ErrGetMethodNotAllowed, o)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func authorizeClient(next http.Handler, o ServerOptions) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("API-Key")
		if key == "" {
			key = r.URL.Query().Get("key")
		}

		if key != o.APIKey {
			ErrorReply(r, w, ErrInvalidAPIKey, o)
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

func setCacheHeaders(next http.Handler, ttl int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer next.ServeHTTP(w, r)

		if r.Method != http.MethodGet || isPublicPath(r.URL.Path) {
			return
		}

		ttlDiff := time.Duration(ttl) * time.Second
		expires := time.Now().Add(ttlDiff)

		w.Header().Add("Expires", strings.ReplaceAll(expires.Format(time.RFC1123), "UTC", "GMT"))
		w.Header().Add("Cache-Control", getCacheControl(ttl))
	})
}

func getCacheControl(ttl int) string {
	if ttl == 0 {
		return "private, no-cache, no-store, must-revalidate"
	}
	return fmt.Sprintf("public, s-maxage=%d, max-age=%d, no-transform", ttl, ttl)
}

func isPublicPath(path string) bool {
	return path == "/" || path == "/health" || path == "/form"
}

func validateURLSignature(next http.Handler, o ServerOptions) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve and remove URL signature from request parameters
		query := r.URL.Query()
		sign := query.Get("sign")
		query.Del("sign")

		// Compute expected URL signature
		h := hmac.New(sha256.New, []byte(o.URLSignatureKey))
		_, _ = h.Write([]byte(r.URL.Path))
		_, _ = h.Write([]byte(query.Encode()))
		expectedSign := h.Sum(nil)

		urlSign, err := base64.RawURLEncoding.DecodeString(sign)
		if err != nil {
			ErrorReply(r, w, ErrInvalidURLSignature, o)
			return
		}

		if !hmac.Equal(urlSign, expectedSign) {
			ErrorReply(r, w, ErrURLSignatureMismatch, o)
			return
		}

		next.ServeHTTP(w, r)
	})
}
