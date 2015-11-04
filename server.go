package main

import (
	"net/http"
	"os"
	"strconv"
	"time"
)

type ServerOptions struct {
	Port             int
	Burst            int
	Concurrency      int
	HttpCacheTtl     int
	HttpReadTimeout  int
	HttpWriteTimeout int
	CORS             bool
	Gzip             bool
	EnableHTTPSource bool
	Address          string
	ApiKey           string
	Mount            string
	CertFile         string
	KeyFile          string
}

func Server(o ServerOptions) error {
	addr := o.Address + ":" + strconv.Itoa(o.Port)
	handler := NewLog(NewServerMux(o), os.Stdout)

	server := &http.Server{
		Addr:           addr,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    time.Duration(o.HttpReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(o.HttpWriteTimeout) * time.Second,
	}

	return listenAndServe(server, o)
}

func listenAndServe(s *http.Server, o ServerOptions) error {
	if o.CertFile != "" && o.KeyFile != "" {
		return s.ListenAndServeTLS(o.CertFile, o.KeyFile)
	}
	return s.ListenAndServe()
}

func NewServerMux(o ServerOptions) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", Middleware(indexController, o))
	mux.Handle("/form", Middleware(formController, o))
	mux.Handle("/health", Middleware(healthController, o))

	image := ImageMiddleware(o)
	mux.Handle("/resize", image(Resize))
	mux.Handle("/enlarge", image(Enlarge))
	mux.Handle("/extract", image(Extract))
	mux.Handle("/crop", image(Crop))
	mux.Handle("/rotate", image(Rotate))
	mux.Handle("/flip", image(Flip))
	mux.Handle("/flop", image(Flop))
	mux.Handle("/thumbnail", image(Thumbnail))
	mux.Handle("/zoom", image(Zoom))
	mux.Handle("/convert", image(Convert))
	mux.Handle("/watermark", image(Watermark))
	mux.Handle("/info", image(Info))

	return mux
}
