# imaginary [![Build Status](https://travis-ci.org/h2non/imaginary.png)](https://travis-ci.org/h2non/imaginary) [![GitHub release](https://img.shields.io/github/tag/h2non/imaginary.svg)](https://github.com/h2non/imaginary/releases) [![GoDoc](https://godoc.org/github.com/h2non/imaginary?status.png)](https://godoc.org/github.com/h2non/imaginary) [![Coverage Status](https://coveralls.io/repos/h2non/imaginary/badge.svg?branch=master)](https://coveralls.io/r/h2non/imaginary?branch=master)

Simple and fast HTTP microservice for image processing powered by [bimg](https://github.com/h2non/bimg) and [libvips](https://github.com/jcupitt/libvips)

`Work in progress`

## Prerequisites

- [libvips](https://github.com/jcupitt/libvips) v7.40.0+ (7.42.0+ recommended)
- C compatible compiler such as gcc 4.6+ or clang 3.0+
- Go 1.3+

## Installation

```bash
go get -u gopkg.in/h2non/imaginary.v0
```

### libvips

Run the following script as `sudo` (supports OSX, Debian/Ubuntu, Redhat, Fedora, Amazon Linux):
```bash
curl -s https://raw.githubusercontent.com/lovell/sharp/master/preinstall.sh | sudo bash -
```

The [install script](https://github.com/lovell/sharp/blob/master/preinstall.sh) requires `curl` and `pkg-config`

## Usage

```
imaginary -p 8080
```

## HTTP API

#### GET /crop

#### POST /crop

## License

MIT - Tomas Aparicio
