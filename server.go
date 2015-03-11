package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Server(addr string, port int) {
	mux := http.NewServeMux()
	mux.Handle("/", middleware(mainHandler))
	mux.Handle("/resize", middleware(processImage))
	//mux.Handle("/crop", middleware(processImage))

	server := &http.Server{
		Addr:           addr + ":" + strconv.Itoa(port),
		Handler:        mux,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	server.ListenAndServe()
}

func handler(fn func(http.ResponseWriter, *http.Request)) http.Handler {
	return middleware(http.HandlerFunc(fn))
}

func middleware(fn func(http.ResponseWriter, *http.Request)) http.Handler {
	next := http.HandlerFunc(fn)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "imgine "+Version)

		logger(w, r)
		validate(next).ServeHTTP(w, r)
	})
}

func logger(w http.ResponseWriter, r *http.Request) {
	remoteAddr := r.Header.Get("X-Forwarded-For")
	if remoteAddr == "" {
		remoteAddr = r.RemoteAddr
	}

	log.Printf("[%s] %s %q\n", r.Method, remoteAddr, r.URL.String())
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("imgine server " + Version))
}

func validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for a request body
		if r.ContentLength == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func processImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	file, mimeType, err := getPayload(r)

	if err != nil || len(mimeType) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//defer r.Body.Close()

	body, err := Resize(file)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", mimeType)
	w.Write(body)
}

type ImageOptions struct {
	Width   int
	Height  int
	Quality int
}

func validateParams(r *http.Request) (*ImageOptions, error) {
	query := r.URL.Query()
	width, _ := strconv.Atoi(query.Get("width"))
	height, _ := strconv.Atoi(query.Get("height"))
	quality, _ := strconv.Atoi(query.Get("quality"))

	if width == 0 || height == 0 {
		return nil, errors.New("Missing required height and width params")
	}

	if quality == 0 {
		quality = 95
	}

	return &ImageOptions{
		Width:   width,
		Height:  height,
		Quality: quality,
	}, nil
}

func getPayload(r *http.Request) ([]byte, string, error) {
	file := r.Body
	mimeType := getMimeType(r)

	if strings.Contains(mimeType, "form-data") {
		file, _, _ = r.FormFile("file")
		mimeType = "image/jpg"
	}

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, mimeType, err
	}

	return buf, mimeType, err
}

func getMimeType(r *http.Request) string {
	mimeType := r.Header.Get("Content-Type")

	if len(mimeType) == 0 {
		mimeType = inferMimeType(r)
	} else {
		mimeType, _, _ = mime.ParseMediaType(mimeType)
	}
	log.Println("Mime " + mimeType)

	return mimeType
}

// pending, slice buffer
func inferMimeType(r *http.Request) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	filetype := http.DetectContentType(buf.Bytes())

	switch filetype {
	case "image/jpeg", "image/jpg", "image/png":
		return filetype
	default:
		return ""
	}
}
