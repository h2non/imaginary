package main

import (
	"github.com/daaku/go.httpgzip"
	"net/http"
	"strconv"
	"time"
)

const maxMemory int64 = 1024 * 1024 * 64

type ServerOptions struct {
	Port    int
	Address string
}

func Server(o ServerOptions) error {
	mux := http.NewServeMux()

	mux.Handle("/", middleware(indexController))
	mux.Handle("/form", middleware(formController))
	mux.Handle("/resize", imageHandler(Resize))
	mux.Handle("/enlarge", imageHandler(Enlarge))
	mux.Handle("/extract", imageHandler(Extract))
	mux.Handle("/crop", imageHandler(Crop))
	mux.Handle("/rotate", imageHandler(Rotate))
	mux.Handle("/flip", imageHandler(Flip))
	mux.Handle("/flop", imageHandler(Flop))
	mux.Handle("/thumbnail", imageHandler(Thumbnail))
	mux.Handle("/zoom", imageHandler(Zoom))
	mux.Handle("/convert", imageHandler(Convert))
	mux.Handle("/watermark", imageHandler(Watermark))
	mux.Handle("/info", imageHandler(Info))

	addr := o.Address + ":" + strconv.Itoa(o.Port)
	server := &http.Server{
		Addr:           addr,
		Handler:        NewLog(mux),
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return server.ListenAndServe()
}

func imageHandler(fn Operation) http.Handler {
	return middleware(mainController(fn))
}

func middleware(fn func(http.ResponseWriter, *http.Request)) http.Handler {
	next := httpgzip.NewHandler(http.HandlerFunc(fn))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "imaginary "+Version)
		validate(next).ServeHTTP(w, r)
	})
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
