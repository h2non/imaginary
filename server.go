package main

import (
	"net/http"
	"net/url"
	"os"
	"path"
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
	EnableURLSource  bool
	AuthForwarding   bool
	Address          string
	PathPrefix       string
	ApiKey           string
	Mount            string
	CertFile         string
	KeyFile          string
	Authorization    string
	AlloweOrigins    []*url.URL
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

	mux.Handle(path.Join(o.PathPrefix, "/"), Middleware(indexController, o))
	mux.Handle(path.Join(o.PathPrefix, "/form"), Middleware(formController, o))
	mux.Handle(path.Join(o.PathPrefix, "/health"), Middleware(healthController, o))

	image := ImageMiddleware(o)
	mux.Handle(path.Join(o.PathPrefix, "/resize"), image(Resize))
	mux.Handle(path.Join(o.PathPrefix, "/enlarge"), image(Enlarge))
	mux.Handle(path.Join(o.PathPrefix, "/extract"), image(Extract))
	mux.Handle(path.Join(o.PathPrefix, "/crop"), image(Crop))
	mux.Handle(path.Join(o.PathPrefix, "/rotate"), image(Rotate))
	mux.Handle(path.Join(o.PathPrefix, "/flip"), image(Flip))
	mux.Handle(path.Join(o.PathPrefix, "/flop"), image(Flop))
	mux.Handle(path.Join(o.PathPrefix, "/thumbnail"), image(Thumbnail))
	mux.Handle(path.Join(o.PathPrefix, "/zoom"), image(Zoom))
	mux.Handle(path.Join(o.PathPrefix, "/convert"), image(Convert))
	mux.Handle(path.Join(o.PathPrefix, "/watermark"), image(Watermark))
	mux.Handle(path.Join(o.PathPrefix, "/info"), image(Info))

	return mux
}
