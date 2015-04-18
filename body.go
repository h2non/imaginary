package main

import (
	"io/ioutil"
	"net/http"
	"strings"
)

const maxMemory int64 = 1024 * 1024 * 64

func readBody(r *http.Request) ([]byte, error) {
	var err error
	var buf []byte

	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/") {
		err = r.ParseMultipartForm(maxMemory)
		if err != nil {
			return nil, err
		}

		buf, err = readFormPayload(r)
	} else {
		buf, err = ioutil.ReadAll(r.Body)
	}

	return buf, err
}

func readFormPayload(r *http.Request) ([]byte, error) {
	file, _, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if len(buf) == 0 {
		return nil, NewError("Empty payload", BAD_REQUEST)
	}

	return buf, err
}
