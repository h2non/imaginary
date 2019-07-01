package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"runtime"
	d "runtime/debug"
	"strconv"
	"strings"
	"time"

	bimg "gopkg.in/h2non/bimg.v1"
)

var (
	aAddr               = flag.String("a", "", "Bind address")
	aPort               = flag.Int("p", 8088, "Port to listen")
	aVers               = flag.Bool("v", false, "Show version")
	aVersl              = flag.Bool("version", false, "Show version")
	aHelp               = flag.Bool("h", false, "Show help")
	aHelpl              = flag.Bool("help", false, "Show help")
	aPathPrefix         = flag.String("path-prefix", "/", "Url path prefix to listen to")
	aCors               = flag.Bool("cors", false, "Enable CORS support")
	aGzip               = flag.Bool("gzip", false, "Enable gzip compression (deprecated)")
	aAuthForwarding     = flag.Bool("enable-auth-forwarding", false, "Forwards X-Forward-Authorization or Authorization header to the image source server. -enable-url-source flag must be defined. Tip: secure your server from public access to prevent attack vectors")
	aEnableURLSource    = flag.Bool("enable-url-source", false, "Enable remote HTTP URL image source processing")
	aEnablePlaceholder  = flag.Bool("enable-placeholder", false, "Enable image response placeholder to be used in case of error")
	aEnableURLSignature = flag.Bool("enable-url-signature", false, "Enable URL signature (URL-safe Base64-encoded HMAC digest)")
	aURLSignatureKey    = flag.String("url-signature-key", "", "The URL signature key (32 characters minimum)")
	aAllowedOrigins     = flag.String("allowed-origins", "", "Restrict remote image source processing to certain origins (separated by commas). Note: Origins are validated against host *AND* path.")
	aMaxAllowedSize     = flag.Int("max-allowed-size", 0, "Restrict maximum size of http image source (in bytes)")
	aKey                = flag.String("key", "", "Define API key for authorization")
	aMount              = flag.String("mount", "", "Mount server local directory")
	aCertFile           = flag.String("certfile", "", "TLS certificate file path")
	aKeyFile            = flag.String("keyfile", "", "TLS private key file path")
	aAuthorization      = flag.String("authorization", "", "Defines a constant Authorization header value passed to all the image source servers. -enable-url-source flag must be defined. This overwrites authorization headers forwarding behavior via X-Forward-Authorization")
	aForwardHeaders     = flag.String("forward-headers", "", "Forwards custom headers to the image source server. -enable-url-source flag must be defined.")
	aPlaceholder        = flag.String("placeholder", "", "Image path to image custom placeholder to be used in case of error. Recommended minimum image size is: 1200x1200")
	aDisableEndpoints   = flag.String("disable-endpoints", "", "Comma separated endpoints to disable. E.g: form,crop,rotate,health")
	aHTTPCacheTTL       = flag.Int("http-cache-ttl", -1, "The TTL in seconds")
	aReadTimeout        = flag.Int("http-read-timeout", 60, "HTTP read timeout in seconds")
	aWriteTimeout       = flag.Int("http-write-timeout", 60, "HTTP write timeout in seconds")
	aConcurrency        = flag.Int("concurrency", 0, "Throttle concurrency limit per second")
	aBurst              = flag.Int("burst", 100, "Throttle burst max cache size")
	aMRelease           = flag.Int("mrelease", 30, "OS memory release interval in seconds")
	aCpus               = flag.Int("cpus", runtime.GOMAXPROCS(-1), "Number of cpu cores to use")
)

const usage = `imaginary %s

Usage:
  imaginary -p 80
  imaginary -cors
  imaginary -concurrency 10
  imaginary -path-prefix /api/v1
  imaginary -enable-url-source
  imaginary -disable-endpoints form,health,crop,rotate
  imaginary -enable-url-source -allowed-origins http://localhost,http://server.com
  imaginary -enable-url-source -enable-auth-forwarding
  imaginary -enable-url-source -authorization "Basic AwDJdL2DbwrD=="
  imaginary -enable-placeholder
  imaginary -enable-url-source -placeholder ./placeholder.jpg
  imaginary -enable-url-signature -url-signature-key 4f46feebafc4b5e988f131c4ff8b5997
  imaginary -enable-url-source -forward-headers X-Custom,X-Token
  imaginary -h | -help
  imaginary -v | -version

Options:
  -a <addr>                 Bind address [default: *]
  -p <port>                 Bind port [default: 8088]
  -h, -help                 Show help
  -v, -version              Show version
  -path-prefix <value>      Url path prefix to listen to [default: "/"]
  -cors                     Enable CORS support [default: false]
  -gzip                     Enable gzip compression (deprecated) [default: false]
  -disable-endpoints        Comma separated endpoints to disable. E.g: form,crop,rotate,health [default: ""]
  -key <key>                Define API key for authorization
  -mount <path>             Mount server local directory
  -http-cache-ttl <num>     The TTL in seconds. Adds caching headers to locally served files.
  -http-read-timeout <num>  HTTP read timeout in seconds [default: 30]
  -http-write-timeout <num> HTTP write timeout in seconds [default: 30]
  -enable-url-source        Enable remote HTTP URL image source processing
  -enable-placeholder       Enable image response placeholder to be used in case of error [default: false]
  -enable-auth-forwarding   Forwards X-Forward-Authorization or Authorization header to the image source server. -enable-url-source flag must be defined. Tip: secure your server from public access to prevent attack vectors
  -forward-headers          Forwards custom headers to the image source server. -enable-url-source flag must be defined.
  -enable-url-signature     Enable URL signature (URL-safe Base64-encoded HMAC digest) [default: false]
  -url-signature-key        The URL signature key (32 characters minimum)
  -allowed-origins <urls>   Restrict remote image source processing to certain origins (separated by commas)
  -max-allowed-size <bytes> Restrict maximum size of http image source (in bytes)
  -certfile <path>          TLS certificate file path
  -keyfile <path>           TLS private key file path
  -authorization <value>    Defines a constant Authorization header value passed to all the image source servers. -enable-url-source flag must be defined. This overwrites authorization headers forwarding behavior via X-Forward-Authorization
  -placeholder <path>       Image path to image custom placeholder to be used in case of error. Recommended minimum image size is: 1200x1200
  -concurrency <num>        Throttle concurrency limit per second [default: disabled]
  -burst <num>              Throttle burst max cache size [default: 100]
  -mrelease <num>           OS memory release interval in seconds [default: 30]
  -cpus <num>               Number of used cpu cores.
                            (default for current machine is %d cores)
`

type URLSignature struct {
	Key string
}

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
	urlSignature := getURLSignature(*aURLSignatureKey)

	opts := ServerOptions{
		Port:               port,
		Address:            *aAddr,
		CORS:               *aCors,
		AuthForwarding:     *aAuthForwarding,
		EnableURLSource:    *aEnableURLSource,
		EnablePlaceholder:  *aEnablePlaceholder,
		EnableURLSignature: *aEnableURLSignature,
		URLSignatureKey:    urlSignature.Key,
		PathPrefix:         *aPathPrefix,
		APIKey:             *aKey,
		Concurrency:        *aConcurrency,
		Burst:              *aBurst,
		Mount:              *aMount,
		CertFile:           *aCertFile,
		KeyFile:            *aKeyFile,
		Placeholder:        *aPlaceholder,
		HTTPCacheTTL:       *aHTTPCacheTTL,
		HTTPReadTimeout:    *aReadTimeout,
		HTTPWriteTimeout:   *aWriteTimeout,
		Authorization:      *aAuthorization,
		ForwardHeaders:     parseForwardHeaders(*aForwardHeaders),
		AllowedOrigins:     parseOrigins(*aAllowedOrigins),
		MaxAllowedSize:     *aMaxAllowedSize,
	}

	// Show warning if gzip flag is passed
	if *aGzip {
		fmt.Println("warning: -gzip flag is deprecated and will not have effect")
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
	if *aHTTPCacheTTL != -1 {
		checkHTTPCacheTTL(*aHTTPCacheTTL)
	}

	// Parse endpoint names to disabled, if present
	if *aDisableEndpoints != "" {
		opts.Endpoints = parseEndpoints(*aDisableEndpoints)
	}

	// Read placeholder image, if required
	if *aPlaceholder != "" {
		buf, err := ioutil.ReadFile(*aPlaceholder)
		if err != nil {
			exitWithError("cannot start the server: %s", err)
		}

		imageType := bimg.DetermineImageType(buf)
		if !bimg.IsImageTypeSupportedByVips(imageType).Load {
			exitWithError("Placeholder image type is not supported. Only JPEG, PNG or WEBP are supported")
		}

		opts.PlaceholderImage = buf
	} else if *aEnablePlaceholder {
		// Expose default placeholder
		opts.PlaceholderImage = placeholder
	}

	// Check URL signature key, if required
	if *aEnableURLSignature {
		if urlSignature.Key == "" {
			exitWithError("URL signature key is required")
		}

		if len(urlSignature.Key) < 32 {
			exitWithError("URL signature key must be a minimum of 32 characters")
		}
	}

	debug("imaginary server listening on port :%d/%s", opts.Port, strings.TrimPrefix(opts.PathPrefix, "/"))

	// Load image source providers
	LoadSources(opts)

	// Start the server
	err := Server(opts)
	if err != nil {
		exitWithError("cannot start the server: %s", err)
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

func getURLSignature(key string) URLSignature {
	if keyEnv := os.Getenv("URL_SIGNATURE_KEY"); keyEnv != "" {
		key = keyEnv
	}

	return URLSignature{key}
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
		exitWithError("error while mounting directory: %s", err)
	}
	if !src.IsDir() {
		exitWithError("mount path is not a directory: %s", path)
	}
	if path == "/" {
		exitWithError("cannot mount root directory for security reasons")
	}
}

func checkHTTPCacheTTL(ttl int) {
	if ttl < 0 || ttl > 31556926 {
		exitWithError("The -http-cache-ttl flag only accepts a value from 0 to 31556926")
	}

	if ttl == 0 {
		debug("Adding HTTP cache control headers set to prevent caching.")
	}
}

func parseForwardHeaders(forwardHeaders string) []string {
	var headers []string
	if forwardHeaders == "" {
		return headers
	}

	for _, header := range strings.Split(forwardHeaders, ",") {
		if norm := strings.TrimSpace(header); norm != "" {
			headers = append(headers, norm)
		}
	}
	return headers
}

func parseOrigins(origins string) []*url.URL {
	var urls []*url.URL
	if origins == "" {
		return urls
	}
	for _, origin := range strings.Split(origins, ",") {
		u, err := url.Parse(origin)
		if err != nil {
			continue
		}

		if u.Path != "" && u.Path[len(u.Path)-1:] != "/" {
			u.Path += "/"
		}

		urls = append(urls, u)
	}
	return urls
}

func parseEndpoints(input string) Endpoints {
	var endpoints Endpoints
	for _, endpoint := range strings.Split(input, ",") {
		endpoint = strings.ToLower(strings.TrimSpace(endpoint))
		if endpoint != "" {
			endpoints = append(endpoints, endpoint)
		}
	}
	return endpoints
}

func memoryRelease(interval int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	go func() {
		for range ticker.C {
			debug("FreeOSMemory()")
			d.FreeOSMemory()
		}
	}()
}

func exitWithError(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", args)
	os.Exit(1)
}

func debug(msg string, values ...interface{}) {
	debug := os.Getenv("DEBUG")
	if debug == "imaginary" || debug == "*" {
		log.Printf(msg, values...)
	}
}
