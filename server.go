package main

import (
	"errors"
	"github.com/daaku/go.httpgzip"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const maxMemory int64 = 1024 * 1024 * 1024

func NewServer(port int) {
	mux := http.NewServeMux()
	mux.Handle("/", middleware(indexHandler))
	mux.Handle("/resize", middleware(processImage))
	mux.Handle("/crop", middleware(processImage))

	server := &http.Server{
		Addr:           ":" + strconv.Itoa(port),
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
	next := httpgzip.NewHandler(http.HandlerFunc(fn))
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

func validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for a request body
		if r.Method != "GET" && r.ContentLength == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("imgine server " + Version))
}

func processImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		contentType = http.DetectContentType(buf)
		contentType = strings.Split(contentType, ";")[0]
	}

	// temporal
	if !strings.HasPrefix(contentType, "multipart/") {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	err := r.ParseMultipartForm(maxMemory)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	file, mimeType, err := getPayload(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if mimeType != "image/jpeg" && mimeType != "image/png" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	body, err := Resize(file)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", mimeType)
	w.Write(body)
}

func getPayload(r *http.Request) ([]byte, string, error) {
	file, _, err := r.FormFile("file")
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, "", err
	}
	if len(buf) == 0 {
		return nil, "", errors.New("Empty payload")
	}

	mimeType := http.DetectContentType(buf)

	return buf, mimeType, err
}
