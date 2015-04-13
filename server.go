package main

import (
	"net/http"
	"os"
	//"runtime"
	"strconv"
	"time"
	//"github.com/PuerkitoBio/throttled"
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
	handler := NewLog(mux, os.Stdout)

	// Throttle by interval
	//th := throttled.Interval(throttled.PerSec(10), 100, &throttled.VaryBy{Path: true}, 50)
	//h := th.Throttle(myHandler)
	//http.ListenAndServe(":9000", h)
	// Throttle by memory
	//th := throttled.MemStats(throttled.MemThresholds(&runtime.MemStats{NumGC: 10}, 10*time.Millisecond)
	//h := th.Throttle(myHandler)
	//http.ListenAndServe(":9000", h)
	// Throttle by rate
	//th := throttled.RateLimit(throttled.PerMin(30), &throttled.VaryBy{RemoteAddr: true}, store.NewMemStore(1000))
	//h := th.Throttle(myHandler)
	//http.ListenAndServe(":9000", h)

	server := &http.Server{
		Addr:           addr,
		Handler:        handler,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return server.ListenAndServe()
}
