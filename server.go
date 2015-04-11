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

const maxMemory int64 = 1024 * 1024 * 64

type ServerOptions struct {
	Port    int
	Address string
}

func NewServer(o ServerOptions) error {
	mux := http.NewServeMux()
	mux.Handle("/form", middleware(uploadForm))
	mux.Handle("/extract", middleware(processImage))
	mux.Handle("/enlarge", middleware(processImage))
	mux.Handle("/resize", middleware(processImage))
	mux.Handle("/crop", middleware(processImage))
	mux.Handle("/thumbnail", middleware(processImage))
	mux.Handle("/rotate", middleware(processImage))
	mux.Handle("/flip", middleware(processImage))
	mux.Handle("/flop", middleware(processImage))
	mux.Handle("/zoom", middleware(processImage))
	mux.Handle("/format", middleware(processImage))
	mux.Handle("/convert", middleware(processImage))
	mux.Handle("/watermark", middleware(processImage))
	mux.Handle("/", middleware(indexHandler))

	addr := o.Address + ":" + strconv.Itoa(o.Port)
	server := &http.Server{
		Addr:           addr,
		Handler:        mux,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return server.ListenAndServe()
}

func handler(fn func(http.ResponseWriter, *http.Request)) http.Handler {
	return middleware(http.HandlerFunc(fn))
}

func middleware(fn func(http.ResponseWriter, *http.Request)) http.Handler {
	next := httpgzip.NewHandler(http.HandlerFunc(fn))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "imaginary "+Version)
		logger(r)
		validateRequest(next).ServeHTTP(w, r)
	})
}

func logger(r *http.Request) {
	remoteAddr := r.Header.Get("X-Forwarded-For")
	if remoteAddr == "" {
		remoteAddr = r.RemoteAddr
	}
	log.Printf("[%s] %s %q\n", r.Method, remoteAddr, r.URL.String())
}

func validateRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" && r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if r.Method == "POST" && r.ContentLength == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("imaginary server " + Version))
}

const formText = `
<html>
<body>
<h1>Resize</h1>
<form method="POST" action="/resize?width=400&height=300" enctype="multipart/form-data">
  <input type="file" name="file" />
  <input type="submit" value="Upload" />
</form>
<h1>Crop</h1>
<form method="POST" action="/crop" enctype="multipart/form-data">
  <input type="file" name="file" />
  <input type="submit" value="Upload" />
</form>
<h1>Flip</h1>
<form method="POST" action="/flip" enctype="multipart/form-data">
  <input type="file" name="file" />
  <input type="submit" value="Upload" />
</form>
<h1>Thumbnail</h1>
<form method="POST" action="/thumbnail" enctype="multipart/form-data">
  <input type="file" name="file" />
  <input type="submit" value="Upload" />
</form>
</body>
</html>
`

func uploadForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(formText))
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
		w.WriteHeader(http.StatusBadRequest)
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
