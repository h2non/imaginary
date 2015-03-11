package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
)

const Usage = `
  Usage:
    imgine [--address a]
    imgine -h | --help
    imgine --version

  Options:
    -a, --address a   bind address [default: *:8088]
    -p, --port 8088   HTTP server port [default: 8088]
    -h, --help        output help information
    -v, --version     output version
`

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	port := flag.Int("port", 8088, "port to listen")
	//addr := flag.String("address", "", "bind address")
	flag.Parse()

	if portEnv := os.Getenv("PORT"); portEnv != "" {
		newPort, _ := strconv.Atoi(portEnv)
		if newPort > 0 {
			*port = newPort
		}
	}

	defer Server("", 8088)

	fmt.Fprintf(os.Stdin, "imgine server listening on port %v\n", *port)
}
