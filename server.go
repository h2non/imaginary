package main

import (
	"net/http"
	"os"
	"strconv"
	"time"
)

type ServerOptions struct {
	Port    int
	CORS    bool
	Gzip    bool
	Address string
	ApiKey  string
}

func Server(o ServerOptions) error {
	mux := http.NewServeMux()

	image := ImageMiddleware(o)
	mux.Handle("/", Middleware(indexController, o))
	mux.Handle("/form", Middleware(formController, o))
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

	addr := o.Address + ":" + strconv.Itoa(o.Port)
	server := &http.Server{
		Addr:           addr,
		Handler:        NewLog(mux, os.Stdout),
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return server.ListenAndServe()
}
