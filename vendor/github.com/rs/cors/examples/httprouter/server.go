package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

func Hello(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"hello\": \"world\"}"))
}

func main() {
	router := httprouter.New()
	router.GET("/", Hello)

	handler := cors.Default().Handler(router)

	http.ListenAndServe(":8080", handler)
}
