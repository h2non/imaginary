package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
)

var (
	p      = flag.Int("p", 8088, "port to listen")
	v      = flag.Bool("v", false, "")
	c      = flag.Int("c", 50, "")
	n      = flag.Int("n", 200, "")
	cpus   = flag.Int("cpus", runtime.GOMAXPROCS(-1), "")
	memory = flag.Int64("memory", maxMemory, "")
	help   = flag.Bool("help", false, "")
)

const usage = `imgine server %s

Usage:
  imgine [--address a]
  imgine -h | --help
  imgine --version

Options:
  -a        bind address [default: *:8088]
  -p        HTTP server port [default: 8088]
  -h        output help
  -v        output version
  -cpus     Number of used cpu cores.
            (default for current machine is %d cores)
  -memory   Max vips memory limit. Defaul to 100MB
`

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(usage, Version, runtime.NumCPU()))
	}
	flag.Parse()

	showHelp := *help
	showVersion := *v
	port := *p

	if showHelp {
		showUsage()
	}

	if showVersion {
		fmt.Fprintln(os.Stdout, Version)
		return
	}

	runtime.GOMAXPROCS(*cpus)

	if portEnv := os.Getenv("PORT"); portEnv != "" {
		newPort, _ := strconv.Atoi(portEnv)
		if newPort > 0 {
			port = newPort
		}
	}

	defer NewServer(port)

	fmt.Fprintf(os.Stdin, "imgine server listening on port %d\n", port)
}

func showUsage() {
	flag.Usage()
	os.Exit(1)
}
