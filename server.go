package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type ServerOptions struct {
	Port               int
	Burst              int
	UnixSocket         string
	Concurrency        int
	HTTPCacheTTL       int
	HTTPReadTimeout    int
	HTTPWriteTimeout   int
	MaxAllowedSize     int
	MaxAllowedPixels   float64
	CORS               bool
	Gzip               bool // deprecated
	AuthForwarding     bool
	EnableURLSource    bool
	EnablePlaceholder  bool
	EnableURLSignature bool
	URLSignatureKey    string
	Address            string
	PathPrefix         string
	APIKey             string
	Mount              string
	CertFile           string
	KeyFile            string
	Authorization      string
	Placeholder        string
	PlaceholderStatus  int
	ForwardHeaders     []string
	PlaceholderImage   []byte
	Endpoints          Endpoints
	AllowedOrigins     []*url.URL
	LogLevel           string
	ReturnSize         bool
}

// Endpoints represents a list of endpoint names to disable.
type Endpoints []string

// IsValid validates if a given HTTP request endpoint is valid or not.
func (e Endpoints) IsValid(r *http.Request) bool {
	parts := strings.Split(r.URL.Path, "/")
	endpoint := parts[len(parts)-1]
	for _, name := range e {
		if endpoint == name {
			return false
		}
	}
	return true
}

func Server(o ServerOptions) {
	addr := o.Address + ":" + strconv.Itoa(o.Port)
	handler := NewLog(NewServerMux(o), os.Stdout, o.LogLevel)

	server := &http.Server{
		Addr:           addr,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    time.Duration(o.HTTPReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(o.HTTPWriteTimeout) * time.Second,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := listenAndServe(server, o); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-done
	log.Print("Graceful shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
}

func listenAndServe(s *http.Server, o ServerOptions) error {
	if o.UnixSocket != "" {
		listener, err := net.Listen("unix", o.UnixSocket)
		if err != nil {
			return err
		}
		return s.Serve(listener)
	}
	if o.CertFile != "" && o.KeyFile != "" {
		return s.ListenAndServeTLS(o.CertFile, o.KeyFile)
	}
	return s.ListenAndServe()
}

func join(o ServerOptions, route string) string {
	return path.Join(o.PathPrefix, route)
}

// NewServerMux creates a new HTTP server route multiplexer.
func NewServerMux(o ServerOptions) http.Handler {
	mux := http.NewServeMux()

	mux.Handle(join(o, "/"), Middleware(indexController(o), o))
	mux.Handle(join(o, "/form"), Middleware(formController(o), o))
	mux.Handle(join(o, "/health"), Middleware(healthController, o))

	image := ImageMiddleware(o)
	mux.Handle(join(o, "/resize"), image(Resize))
	mux.Handle(join(o, "/fit"), image(Fit))
	mux.Handle(join(o, "/enlarge"), image(Enlarge))
	mux.Handle(join(o, "/extract"), image(Extract))
	mux.Handle(join(o, "/crop"), image(Crop))
	mux.Handle(join(o, "/smartcrop"), image(SmartCrop))
	mux.Handle(join(o, "/rotate"), image(Rotate))
	mux.Handle(join(o, "/autorotate"), image(AutoRotate))
	mux.Handle(join(o, "/flip"), image(Flip))
	mux.Handle(join(o, "/flop"), image(Flop))
	mux.Handle(join(o, "/thumbnail"), image(Thumbnail))
	mux.Handle(join(o, "/zoom"), image(Zoom))
	mux.Handle(join(o, "/convert"), image(Convert))
	mux.Handle(join(o, "/watermark"), image(Watermark))
	mux.Handle(join(o, "/watermarkimage"), image(WatermarkImage))
	mux.Handle(join(o, "/info"), image(Info))
	mux.Handle(join(o, "/blur"), image(GaussianBlur))
	mux.Handle(join(o, "/pipeline"), image(Pipeline))

	return mux
}
