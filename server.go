package main

import (
	"fmt"
	//"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"strings"
	"time"
)

func Server(port int) {
	server := &http.Server{
		Addr:           ":8088",
		ReadTimeout:    120 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Register this pat with the default serve mux so that other packages
	// may also be exported. (i.e. /debug/pprof/*)
	//server.Handle("/", ProcessImage)
	http.HandleFunc("/", ProcessImage)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting the server:", err)
	}
}

func ProcessImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "imgine "+Version)

	file, mimeType, err := getPayload(r)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(mimeType) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := Resize(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	w.Write(body)
}

func getPayload(r *http.Request) ([]byte, string, error) {
	file := r.Body
	mimeType := getMimeType(r)

	if strings.Contains(mimeType, "form-data") {
		file, _, _ = r.FormFile("file")
	}

	buf, err := ioutil.ReadAll(file)

	return buf, mimeType, err
}

func getMimeType(r *http.Request) string {
	mimeType := r.Header.Get("Content-Type")

	if len(mimeType) == 0 {
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return ""
		}
		mimeType = inferMimeType(buf)
	} else {
		mimeType, _, _ = mime.ParseMediaType(mimeType)
	}

	return mimeType
}

func inferMimeType(buf []byte) string {
	chunk := make([]byte, 512)
	filetype := http.DetectContentType(chunk)

	switch filetype {
	case "image/jpeg", "image/jpg", "image/png":
		return filetype
	default:
		return ""
	}
}
