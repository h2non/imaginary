package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const formatPattern = "%s - - [%s] \"%s\" %d %d %.4f\n"

// LogRecord implements an Apache-compatible HTTP logging
type LogRecord struct {
	http.ResponseWriter
	status                int
	responseBytes         int64
	ip                    string
	method, uri, protocol string
	time                  time.Time
	elapsedTime           time.Duration
}

// Log writes a log entry in the passed io.Writer stream
func (r *LogRecord) Log(out io.Writer) {
	timeFormat := r.time.Format("02/Jan/2006 15:04:05")
	request := fmt.Sprintf("%s %s %s", r.method, r.uri, r.protocol)
	fmt.Fprintf(out, formatPattern, r.ip, timeFormat, request, r.status, r.responseBytes, r.elapsedTime.Seconds())
}

// Write acts like a proxy passing the given bytes buffer to the ResponseWritter
// and additionally counting the passed amount of bytes for logging usage.
func (r *LogRecord) Write(p []byte) (int, error) {
	written, err := r.ResponseWriter.Write(p)
	r.responseBytes += int64(written)
	return written, err
}

// WriteHeader calls ResponseWriter.WriteHeader() and sets the status code
func (r *LogRecord) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// LogHandler maps the HTTP handler with a custom io.Writer compatible stream
type LogHandler struct {
	handler http.Handler
	io      io.Writer
}

// NewLog creates a new logger
func NewLog(handler http.Handler, io io.Writer) http.Handler {
	return &LogHandler{handler, io}
}

// Implements the required method as standard HTTP handler, serving the request.
func (h *LogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	if colon := strings.LastIndex(clientIP, ":"); colon != -1 {
		clientIP = clientIP[:colon]
	}

	record := &LogRecord{
		ResponseWriter: w,
		ip:             clientIP,
		time:           time.Time{},
		method:         r.Method,
		uri:            r.RequestURI,
		protocol:       r.Proto,
		status:         http.StatusOK,
		elapsedTime:    time.Duration(0),
	}

	startTime := time.Now()
	h.handler.ServeHTTP(record, r)
	finishTime := time.Now()

	record.time = finishTime.UTC()
	record.elapsedTime = finishTime.Sub(startTime)

	record.Log(h.io)
}
