# imaginary [![Build Status](https://travis-ci.org/h2non/imaginary.png)](https://travis-ci.org/h2non/imaginary) [![Docker](https://img.shields.io/badge/docker-h2non/imaginary-blue.svg)](https://registry.hub.docker.com/u/h2non/imaginary/) [![Docker Registry](https://img.shields.io/docker/pulls/h2non/imaginary.svg)](https://registry.hub.docker.com/u/h2non/imaginary) [![Heroku](https://img.shields.io/badge/Heroku-Deploy_Now-blue.svg)](https://heroku.com/deploy) [![Go Report Card](http://goreportcard.com/badge/h2non/imaginary)](http://goreportcard.com/report/h2non/imaginary)

<img src="http://s14.postimg.org/8th71a201/imaginary_world.jpg" width="100%" />

**[Fast](#benchmarks) HTTP [microservice](http://microservices.io/patterns/microservices.html)** written in Go **for high-level image processing** backed by [bimg](https://github.com/h2non/bimg) and [libvips](https://github.com/jcupitt/libvips). `imaginary` can be used as private or public HTTP service for massive image processing. 
It's almost dependency-free and only uses [`net/http`](http://golang.org/pkg/net/http/) native package for better [performance](#performance).

Supports multiple [image operations](#supported-image-operations) exposed as a simple [HTTP API](#http-api), 
with additional optional features such as **API token authorization**, **gzip compression**, **HTTP traffic throttle** strategy and **CORS support** for web clients.

`imaginary` **can read** images **from HTTP payloads**, **server local path** or **remote HTTP servers**, supporting **JPEG**, **PNG**, **WEBP** and **TIFF** formats and it's able to output to JPEG, PNG and WEBP, including conversion between them.

It uses internally libvips, a powerful and efficient library written in C for image processing 
which requires a [low memory footprint](http://www.vips.ecs.soton.ac.uk/index.php?title=Speed_and_Memory_Use) 
and it's typically 4x faster than using the quickest ImageMagick and GraphicsMagick 
settings or Go native `image` package, and in some cases it's even 8x faster processing JPEG images. 

To get started, take a look the [installation](#installation) steps, [usage](#usage) cases and [API](#http-api) docs.

`imaginary` is currently used in production processing thousands of images per day.

## Contents

- [Supported image operations](#supported-image-operations)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
  - [Docker](#docker)
  - [Heroku](#heroku)
- [Recommended resources](#recommended-resources)
- [Production notes](#production-notes)
- [Scalability](#scalability)
- [Clients](#clients)
- [Performance](#performance)
- [Benchmark](#benchmark)
- [Usage](#usage)
- [HTTP API](#http-api)
  - [Authorization](#authorization)
  - [Errors](#errors)
  - [Form data](#form-data)
  - [Params](#params)
  - [Endpoints](#get-)

## Supported image operations

- Resize
- Enlarge
- Crop
- Rotate (with auto-rotate based on EXIF orientation)
- Flip (with auto-flip based on EXIF metadata)
- Flop
- Zoom
- Thumbnail
- Extract area
- Watermark (customizable by text)
- Custom output color space (RGB, black/white...)
- Format conversion (with additional quality/compression settings)
- Info (image size, format, orientation, alpha...)

## Prerequisites

- [libvips](https://github.com/jcupitt/libvips) v7.40.0+ (7.42.0+ recommended)
- C compatible compiler such as gcc 4.6+ or clang 3.0+
- Go 1.3+

## Installation

```bash
go get -u github.com/h2non/imaginary
```

### libvips

Run the following script as `sudo` (supports OSX, Debian/Ubuntu, Redhat, Fedora, Amazon Linux):
```bash
curl -s https://raw.githubusercontent.com/lovell/sharp/master/preinstall.sh | sudo bash -
```

The [install script](https://github.com/lovell/sharp/blob/master/preinstall.sh) requires `curl` and `pkg-config`

### Docker

See [Dockerfile](https://github.com/h2non/imaginary/blob/master/Dockerfile) for image details.

Fetch the image (comes with latest stable Go and libvips versions)
```
docker pull h2non/imaginary
```

Start the container with optional flags (default listening on port 9000)
```
docker run -p 9000:9000 h2non/imaginary -cors -gzip
```

Start the container in debug mode:
```
docker run -p 9000:9000 -e "DEBUG=*" h2non/imaginary 
```

Enter to the interactive shell in a running container
```
sudo docker exec -it <containerIdOrName> bash
```

Stop the container
```
docker stop h2non/imaginary
```

You can see all the Docker tags [here](https://hub.docker.com/r/h2non/imaginary/tags/).

### Heroku

Click on the Heroku button to easily deploy your app:

[![Heroku](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy) 

Or alternatively you can follow the manual steps:

Clone this repository:
```
git clone https://github.com/h2non/imaginary.git
```

Set the buildpack for your application
```
heroku config:add BUILDPACK_URL=https://github.com/h2non/heroku-buildpack-imaginary.git
```

Optionally, define the PKGCONFIG path:
```
heroku config:add PKG_CONFIG_PATH=/app/vendor/vips/lib/pkgconfig
```

Add Heroku git remote:
```
heroku git:remote -a your-application
```

Deploy it!
```
git push heroku master
```

### Recommended resources

Given the multithreaded native nature of Go, in term of CPUs, most cores means more concurrency and therefore, a better performance can be achieved. 
From the other hand, in terms of memory, 512MB of RAM is usually enough for small services with low concurrency (<5 requests/second). 
Up to 2GB for high-load HTTP service processing potentially large images or exposed to an eventual high concurrency.

If you need to expose `imaginary` as public HTTP server, it's highly recommended to protect the service against DDoS-like attacks. 
`imaginary` has built-in support for HTTP concurrency throttle strategy to deal with this in a more convenient way and mitigate possible issues limiting the number of concurrent requests per second and caching the awaiting requests, if necessary.

### Production notes

In production focused environments it's highly recommended to enable the HTTP concurrency throttle strategy in your `imaginary` servers.

The recommended concurrency limit per server to guarantee a good performance is up to `20` requests per second.

You can enable it simply passing a flag to the binary:
```
$ imaginary -concurrency 20
```

### Scalability

If you're looking for a large scale solution for massive image processing, you should scale `imaginary` horizontally, distributing the HTTP load across a pool of imaginary servers.

Assuming that you want to provide a high availability to deal efficiently with, let's say, 100 concurrent req/sec, a good approach would be using a front end balancer (e.g: HAProxy) to delegate the traffic control flow, ensure the quality of service and distribution the HTTP across a pool of servers:

```
        |==============|
        |  Dark World  |
        |==============|
              ||||
        |==============|
        |   Balancer   |
        |==============|
           |       |   
          /         \
         /           \
        /             \
 /-----------\   /-----------\
 | imaginary |   | imaginary | (*n)
 \-----------/   \-----------/
```

## Clients

- [node.js/io.js](https://github.com/h2non/node-imaginary)

Feel free to send a PR if you created a client for other language.

## Performance

libvips is probably the faster open source solution for image processing. 
Here you can see some performance test comparisons for multiple scenarios:

- [libvips speed and memory usage](http://www.vips.ecs.soton.ac.uk/index.php?title=Speed_and_Memory_Use)
- [sharp performance tests](https://github.com/lovell/sharp#the-task) 
- [bimg](https://github.com/h2non/bimg#Performance) (Go library with C bindings to libvips)

## Benchmark

See [benchmark.sh](https://github.com/h2non/imaginary/blob/master/benchmark.sh) for more details

Environment: Go 1.4.2. libvips-7.42.3. OSX i7 2.7Ghz

```
Requests  [total]       200
Duration  [total, attack, wait]   10.030639787s, 9.949499515s, 81.140272ms
Latencies [mean, 50, 95, 99, max]   83.124471ms, 82.899435ms, 88.948008ms, 95.547765ms, 104.384977ms
Bytes In  [total, mean]     23443800, 117219.00
Bytes Out [total, mean]     175517000, 877585.00
Success   [ratio]       100.00%
Status Codes  [code:count]      200:200
```

### Conclusions

`imaginary` can deal efficiently with up to 20 request per second running in a multicore machine, 
where it crops a JPEG image of 5MB and spending per each request less than 100 ms

The most expensive image operation under high concurrency scenarios (> 20 req/sec) is the image enlargement, which requires a considerable amount of math operations to scale the original image. In this kind of operation the required processing time usually grows over the time if you're stressing the server continuously. The advice here is as simple as taking care about the number of concurrent enlarge operations to avoid server performance bottlenecks.

## Usage

```
imaginary server

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
  -enable-url-source        Restrict remote image source processing to certain origins (separated by commas)
  -allowed-origins <urls>   TLS certificate file path
  -certfile <path>          TLS certificate file path
  -keyfile <path>           TLS private key file path
  -concurreny <num>         Throttle concurrency limit per second [default: disabled]
  -burst <num>              Throttle burst max cache size [default: 100]
  -mrelease <num>           OS memory release inverval in seconds [default: 30]
  -cpus <num>               Number of used cpu cores.
                            (default for current machine is 8 cores)
```

Start the server in a custom port
```bash
imaginary -p 8080
```

Also, you can pass the port as environment variable
```bash
PORT=8080 imaginary 
```

Enable HTTP server throttle strategy (max 10 requests/second)
```
imaginary -p 8080 -concurrency 10
```

Enable remote URL image fetching (then you can do GET request passing the `url=http://server.com/image.jpg` query param)
```
imaginary -p 8080 -enable-url-source
```

Mount local directory (then you can do GET request passing the `file=image.jpg` query param)
```
imaginary -p 8080 -mount ~/images
```

Send caching headers (only possible with the -mount option). The headers can be set in either "cache nothing" or 
"cache for N seconds". By specifying 0 Imaginary will send the "don't cache" headers, otherwise it sends headers with a 
TTL. The following example informs the client to cache the result for 1 year.
```
imaginary -mount ~/images -http-cache-ttl 31556926
```

Increase libvips threads concurrency (experimental)
```
VIPS_CONCURRENCY=10 imaginary -p 8080 -concurrency 10
```

Enable debug mode
```
DEBUG=* imaginary -p 8080
```

Or filter debug output by package
```
DEBUG=imaginary imaginary -p 8080
```

#### Examples

Reading a local image (you must pass the `-mount=<directory>` flag):
```
curl -O "http://localhost:8088/crop?width=500&height=400&file=foo/bar/image.jpg"
```

Fetching the image from a remote server (you must pass the `-enable-url-source` flag):
```
curl -O "http://localhost:8088/crop?width=500&height=400&url=https://raw.githubusercontent.com/h2non/imaginary/master/fixtures/large.jpg"
```

#### Playground

`imaginary` exposes an ugly HTML form for playground purposes in: [`http://localhost:8088/form`](http://localhost:8088/form)

## HTTP API

### Authorization

imaginary supports a simple token-based API authorization. 
To enable it, you should pass the `-key` flag to the binary.

API token can be defined as HTTP header (`API-Key`) or query param (`key`).

Example request with API key:
```
POST /crop HTTP/1.1
Host: localhost:8088
API-Key: secret
```

### Errors

`imaginary` will always reply with the proper HTTP status code and JSON body with error details.

Here an example response error when the payload is empty:
```json
{
  "message": "Cannot read payload: no such file",
  "code": 1
}
```

See all the predefined supported errors [here](https://github.com/h2non/imaginary/blob/master/error.go#L19-L28).

### Form data

If you're pushing images to `imaginary` as `multipart/form-data` (you can do it as well as `image/*`), you must define at least one input field called `file` with the raw image data in order to be processed properly by imaginary.

### Params

Complete list of available params. Take a look to each specific endpoint to see which params are supported. 
Image measures are always in pixels, unless otherwise indicated.

- **width**       `int`   - Width of image area to extract/resize
- **height**      `int`   - Height of image area to extract/resize 
- **top**         `int`   - Top edge of area to extract. Example: `100`
- **left**        `int`   - Left edge of area to extract. Example: `100`
- **areawidth**   `int`   - Height area to extract. Example: `300`
- **areaheight**  `int`   - Width area to extract. Example: `300`
- **quality**     `int`   - JPEG image quality between 1-100. Defaults to `80`
- **compression** `int`   - PNG compression level. Default: `6`
- **rotate**      `int`   - Image rotation angle. Must be multiple of `90`. Example: `180`
- **factor**      `int`   - Zoom factor level. Example: `2`
- **margin**      `int`   - Text area margin for watermark. Example: `50`
- **dpi**         `int`   - DPI value for watermark. Example: `150`
- **textwidth**   `int`   - Text area width for watermark. Example: `200`
- **opacity**     `float` - Opacity level for watermark text. Default: `0.2`
- **force**       `bool`  - Force image transformation size. Default: `false`
- **nocrop**      `bool`  - Disable crop transformation enabled by default by some operations. Default: `false`
- **noreplicate** `bool`  - Disable text replication in watermark. Defaults to `false`
- **norotation**  `bool`  - Disable auto rotation based on EXIF orientation. Defaults to `false`
- **noprofile**   `bool`  - Disable adding ICC profile metadata. Defaults to `false`
- **text**        `string` - Watermark text content. Example: `copyright (c) 2189`
- **font**        `string` - Watermark text font type and format. Example: `sans bold 12`
- **color**       `string` - Watermark text RGB decimal base color. Example: `255,200,150`
- **type**        `string` - Specify the image format to output. Possible values are: `jpeg`, `png` and `webp`
- **gravity**     `string` - Define the crop operation gravity. Supported values are: `north`, `south`, `centre`, `west` and `east`. Defaults to `centre`.
- **file**        `string` - Use image from server local file path. In order to use this you must pass the `-mount=<dir>` flag.
- **url**         `string` - Fetch the image from a remove HTTP server. In order to use this you must pass the `-enable-url-source` flag.
- **colorspace**  `string` - Use a custom color space for the output image. Allowed values are: `srgb` or `bw` (black&white)
- **field**       `string` - Custom image form field name if using `multipart/form`. Defaults to: `file`

#### GET /
Content-Type: `application/json`

Serves as JSON the current imaginary, bimg and libvips versions.

#### GET /health
Content-Type: `application/json`

Provides some useful statistics about the server stats with the following structure:

- **uptime** `number` - Server process uptime in seconds.
- **allocatedMemory** `number` - Currently allocated memory in megabytes.
- **totalAllocatedMemory** `number` - Total allocated memory over the time in megabytes.
- **gorouting** `number` - Number of running gorouting.
- **cpus** `number` - Number of used CPU cores.

Example response:
```json
{
  "uptime": 1293,
  "allocatedMemory": 5.31,
  "totalAllocatedMemory": 34.3,
  "goroutines": 19,
  "cpus": 8
}
```

#### GET /form
Content Type: `text/html`

Serves an ugly HTML form, just for testing/playground purposes

#### GET | POST /info
Accepts: `image/*, multipart/form-data`. Content-Type: `application/json` 

Returns the image metadata as JSON:
```json
{
  "width": 550,
  "height": 740,
  "type": "jpeg",
  "space": "srgb",
  "hasAlpha": false,
  "hasProfile": true,
  "channels": 3,
  "orientation": 1
}
```

#### GET | POST /crop
Accepts: `image/*, multipart/form-data`. Content-Type: `image/*` 

Crop the image by a given width or height. Image ratio is maintained

##### Allowed params

- width `int`
- height `int`
- quality `int` (JPEG-only)
- compression `int` (PNG-only)
- type `string`
- file `string` - Only GET method and if the `-mount` flag is present
- url `string` - Only GET method and if the `-enable-url-source` flag is present
- force `bool`
- rotate `int`
- norotation `bool`
- noprofile `bool`
- colorspace `string`
- gravity `string`
- field `string` - Only POST and `multipart/form` payloads

#### GET | POST /resize
Accepts: `image/*, multipart/form-data`. Content-Type: `image/*` 

Resize an image by width or height. Image aspect ratio is maintained

##### Allowed params

- width `int` `required`
- height `int`
- quality `int` (JPEG-only)
- compression `int` (PNG-only)
- type `string`
- file `string` - Only GET method and if the `-mount` flag is present
- url `string` - Only GET method and if the `-enable-url-source` flag is present
- force `bool`
- rotate `int`
- norotation `bool`
- noprofile `bool`
- colorspace `string`
- field `string` - Only POST and `multipart/form` payloads

#### GET | POST /enlarge
Accepts: `image/*, multipart/form-data`. Content-Type: `image/*` 

##### Allowed params

- width `int` `required`
- height `int` `required`
- quality `int` (JPEG-only)
- compression `int` (PNG-only)
- type `string`
- file `string` - Only GET method and if the `-mount` flag is present
- url `string` - Only GET method and if the `-enable-url-source` flag is present
- force `bool`
- rotate `int`
- norotation `bool`
- noprofile `bool`
- colorspace `string`
- field `string` - Only POST and `multipart/form` payloads

#### GET | POST /extract
Accepts: `image/*, multipart/form-data`. Content-Type: `image/*` 

##### Allowed params

- top `int` `required`
- left `int`
- areawidth `int` `required`
- areaheight `int`
- width `int`
- height `int` 
- quality `int` (JPEG-only)
- compression `int` (PNG-only)
- type `string`
- file `string` - Only GET method and if the `-mount` flag is present
- url `string` - Only GET method and if the `-enable-url-source` flag is present
- force `bool`
- rotate `int`
- norotation `bool`
- noprofile `bool`
- colorspace `string`
- field `string` - Only POST and `multipart/form` payloads

#### GET | POST /zoom
Accepts: `image/*, multipart/form-data`. Content-Type: `image/*` 

##### Allowed params

- factor `number` `required`
- width `int`
- height `int` 
- quality `int` (JPEG-only)
- compression `int` (PNG-only)
- type `string` 
- file `string` - Only GET method and if the `-mount` flag is present
- url `string` - Only GET method and if the `-enable-url-source` flag is present
- force `bool`
- rotate `int`
- norotation `bool`
- noprofile `bool`
- colorspace `string`
- field `string` - Only POST and `multipart/form` payloads

#### GET | POST /thumbnail
Accepts: `image/*, multipart/form-data`. Content-Type: `image/*` 

##### Allowed params

- width `int`
- height `int` 
- quality `int` (JPEG-only)
- compression `int` (PNG-only)
- type `string`
- file `string` - Only GET method and if the `-mount` flag is present
- url `string` - Only GET method and if the `-enable-url-source` flag is present
- force `bool`
- rotate `int`
- norotation `bool`
- noprofile `bool`
- colorspace `string`
- field `string` - Only POST and `multipart/form` payloads

#### GET | POST /rotate
Accepts: `image/*, multipart/form-data`. Content-Type: `image/*` 

##### Allowed params

- rotate `int` `required`
- width `int`
- height `int` 
- quality `int` (JPEG-only)
- compression `int` (PNG-only)
- type `string`
- file `string` - Only GET method and if the `-mount` flag is present
- url `string` - Only GET method and if the `-enable-url-source` flag is present
- force `bool`
- norotation `bool`
- noprofile `bool`
- colorspace `string`
- field `string` - Only POST and `multipart/form` payloads

#### GET | POST /flip
Accepts: `image/*, multipart/form-data`. Content-Type: `image/*` 

##### Allowed params

- width `int`
- height `int` 
- quality `int` (JPEG-only)
- compression `int` (PNG-only)
- type `string`
- file `string` - Only GET method and if the `-mount` flag is present
- url `string` - Only GET method and if the `-enable-url-source` flag is present
- force `bool`
- norotation `bool`
- noprofile `bool`
- colorspace `string`
- field `string` - Only POST and `multipart/form` payloads

#### GET | POST /flop
Accepts: `image/*, multipart/form-data`. Content-Type: `image/*` 

##### Allowed params

- width `int`
- height `int` 
- quality `int` (JPEG-only)
- compression `int` (PNG-only)
- type `string`
- file `string` - Only GET method and if the `-mount` flag is present
- url `string` - Only GET method and if the `-enable-url-source` flag is present
- force `bool`
- norotation `bool`
- noprofile `bool`
- colorspace `string`
- field `string` - Only POST and `multipart/form` payloads

#### GET | POST /convert
Accepts: `image/*, multipart/form-data`. Content-Type: `image/*` 

##### Allowed params

- type `string` `required`
- quality `int` (JPEG-only)
- compression `int` (PNG-only)
- file `string` - Only GET method and if the `-mount` flag is present
- url `string` - Only GET method and if the `-enable-url-source` flag is present
- force `bool`
- rotate `int`
- norotation `bool`
- noprofile `bool`
- colorspace `string`
- field `string` - Only POST and `multipart/form` payloads

#### GET | POST /watermark
Accepts: `image/*, multipart/form-data`. Content-Type: `image/*` 

##### Allowed params

- text `string` `required`
- margin `int`
- dpi `int`
- textwidth `int`
- opacity `float`
- noreplicate `bool`
- font `string`
- color `string` 
- quality `int` (JPEG-only)
- compression `int` (PNG-only)
- type `string`
- file `string` - Only GET method and if the `-mount` flag is present
- url `string` - Only GET method and if the `-enable-url-source` flag is present
- force `bool`
- rotate `int`
- norotation `bool`
- noprofile `bool`
- colorspace `string`
- field `string` - Only POST and `multipart/form` payloads

## License

MIT - Tomas Aparicio

[![views](https://sourcegraph.com/api/repos/github.com/h2non/imaginary/.counters/views.svg)](https://sourcegraph.com/github.com/h2non/imaginary)
