package main

import (
	"bytes"
	//"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
)

func StartServer(t *testing.T) {
	// start the server
	_, err := Server(8088)
	if err != nil {
		t.Error("Cannot start the server")
	}
}

func TestUploadForm(t *testing.T) {
	//var server *http.Server
	url := "http://localhost:8088"
	file := "fixtures/large.jpg"

	/*
		  defer (func() {
				_, err := Server(8088)
				if err != nil {
					t.Error("Cannot start the server")
				}
			})()

			fmt.Println("Start server")
	*/

	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	// Add your image file
	f, err := os.Open(file)
	if err != nil {
		return
	}
	fw, err := w.CreateFormFile("image", file)
	if err != nil {
		return
	}

	io.Copy(fw, f)
	w.CreateFormField("key")
	fw.Write([]byte("KEY"))

	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", "multipart/form-data")

	// Submit the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Error("Cannot send request ", err.Error())
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		t.Error("Invalid response code ", res.StatusCode)
	}
}
