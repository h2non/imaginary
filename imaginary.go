package main

import (
	"flag"
	"fmt"
	. "github.com/tj/go-debug"
	"net/url"
	"os"
	"runtime"
	d "runtime/debug"
	"strconv"
	"strings"
	"time"
)

var debug = Debug("imaginary")

var (
	aAddr            = flag.String("a", "", "bind address")
	aPort            = flag.Int("p", 8088, "port to listen")
	aVers            = flag.Bool("v", false, "Show version")
	aVersl           = flag.Bool("version", false, "Show version")
	aHelp            = flag.Bool("h", false, "Show help")
	aHelpl           = flag.Bool("help", false, "Show help")
	aCors            = flag.Bool("cors", false, "Enable CORS support")
	aGzip            = flag.Bool("gzip", false, "Enable gzip compression")
	aEnableURLSource = flag.Bool("enable-url-source", false, "Enable remote HTTP URL image source processing")
	aAlloweOrigins   = flag.String("allowed-origins", "", "Restrict remote image source processing to certain origins (separated by commas)")
	aKey             = flag.String("key", "", "Define API key for authorization")
	aMount           = flag.String("mount", "", "Mount server local directory")
	aCertFile        = flag.String("certfile", "", "TLS certificate file path")
	aKeyFile         = flag.String("keyfile", "", "TLS private key file path")
	aHttpCacheTtl    = flag.Int("http-cache-ttl", -1, "The TTL in seconds")
	aReadTimeout     = flag.Int("http-read-timeout", 60, "HTTP read timeout in seconds")
	aWriteTimeout    = flag.Int("http-write-timeout", 60, "HTTP write timeout in seconds")
	aConcurrency     = flag.Int("concurrency", 0, "Throttle concurrency limit per second")
	aBurst           = flag.Int("burst", 100, "Throttle burst max cache size")
	aMRelease        = flag.Int("mrelease", 30, "OS memory release inverval in seconds")
	aCpus            = flag.Int("cpus", runtime.GOMAXPROCS(-1), "Number of cpu cores to use")
)

const usage = `imaginary %s

Usage:
  imaginary -p 80
  imaginary -cors -gzip
  imaginary -concurrency 10
  imaginary -enable-url-source
  imaginary -enable-url-source -allowed-origins http://localhost,http://server.com
  imaginary -h | -help
  imaginary -v | -version

Options:
  -a <addr>                 bind address [default: *]
  -p <port>                 bind port [default: 8088]
  -h, -help                 output help
  -v, -version              output version
  -cors                     Enable CORS support [default: false]
  -gzip                     Enable gzip compression [default: false]
  -key <key>                Define API key for authorization
  -mount <path>             Mount server local directory
  -http-cache-ttl <num>     The TTL in seconds. Adds caching headers to locally served files.
  -http-read-timeout <num>  HTTP read timeout in seconds [default: 30]
  -http-write-timeout <num> HTTP write timeout in seconds [default: 30]
  -enable-url-source        Enable remote HTTP URL image source processing
  -allowed-origins <urls>   Restrict remote image source processing to certain origins (separated by commas)
  -certfile <path>          TLS certificate file path
  -keyfile <path>           TLS private key file path
  -concurreny <num>         Throttle concurrency limit per second [default: disabled]
  -burst <num>              Throttle burst max cache size [default: 100]
  -mrelease <num>           OS memory release inverval in seconds [default: 30]
  -cpus <num>               Number of used cpu cores.
                            (default for current machine is %d cores)
`

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(usage, Version, runtime.NumCPU()))
	}
	flag.Parse()

	if *aHelp || *aHelpl {
		showUsage()
	}
	if *aVers || *aVersl {
		showVersion()
	}

	// Only required in Go < 1.5
	runtime.GOMAXPROCS(*aCpus)

	port := getPort(*aPort)
	opts := ServerOptions{
		Port:             port,
		Address:          *aAddr,
		Gzip:             *aGzip,
		CORS:             *aCors,
		EnableURLSource:  *aEnableURLSource,
		ApiKey:           *aKey,
		Concurrency:      *aConcurrency,
		Burst:            *aBurst,
		Mount:            *aMount,
		CertFile:         *aCertFile,
		KeyFile:          *aKeyFile,
		HttpCacheTtl:     *aHttpCacheTtl,
		HttpReadTimeout:  *aReadTimeout,
		HttpWriteTimeout: *aWriteTimeout,
		AlloweOrigins:    parseOrigins(*aAlloweOrigins),
	}

	// Create a memory release goroutine
	if *aMRelease > 0 {
		memoryRelease(*aMRelease)
	}

	// Check if the mount directory exists, if present
	if *aMount != "" {
		checkMountDirectory(*aMount)
	}

	// Validate HTTP cache param, if present
	if *aHttpCacheTtl != -1 {
		checkHttpCacheTtl(*aHttpCacheTtl)
	}

	debug("imaginary server listening on port %d", port)

	// Load image source providers
	LoadSources(opts)

	// Start the server
	err := Server(opts)
	if err != nil {
		exitWithError("cannot start the server: %s\n", err)
	}
}

func getPort(port int) int {
	if portEnv := os.Getenv("PORT"); portEnv != "" {
		newPort, _ := strconv.Atoi(portEnv)
		if newPort > 0 {
			port = newPort
		}
	}
	return port
}

func showUsage() {
	flag.Usage()
	os.Exit(1)
}

func showVersion() {
	fmt.Println(Version)
	os.Exit(1)
}

func checkMountDirectory(path string) {
	src, err := os.Stat(path)
	if err != nil {
		exitWithError("error while mounting directory: %s\n", err)
	}
	if src.IsDir() == false {
		exitWithError("mount path is not a directory: %s\n", path)
	}
	if path == "/" {
		exitWithError("cannot mount root directory for security reasons")
	}
}

func checkHttpCacheTtl(ttl int) {
	if ttl < -1 || ttl > 31556926 {
		exitWithError("The -http-cache-ttl flag only accepts a value from 0 to 31556926")
	}

	if ttl == 0 {
		debug("Adding HTTP cache control headers set to prevent caching.")
	}
}

func parseOrigins(origins string) []*url.URL {
	urls := []*url.URL{}
	if origins == "" {
		return urls
	}
	for _, origin := range strings.Split(origins, ",") {
		u, err := url.Parse(origin)
		if err != nil {
			continue
		}
		urls = append(urls, u)
	}
	return urls
}

func memoryRelease(interval int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	go func() {
		for _ = range ticker.C {
			debug("FreeOSMemory()")
			d.FreeOSMemory()
		}
	}()
}

func exitWithError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args)
	os.Exit(1)
}
