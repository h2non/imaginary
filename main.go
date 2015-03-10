package main

import (
	"log"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	server, err := Server(8088)
	if err != nil {
		log.Fatalln("Cannot start the server")
	}
}
