# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM marcbachmann/libvips:latest
MAINTAINER tomas@aparicio.me

# Server port to listen
ENV PORT 9000

# Go version to use
ENV GOLANG_VERSION 1.7.1

# gcc for cgo
RUN apt-get update && apt-get install -y \
    gcc curl git libc6-dev make ca-certificates \
    --no-install-recommends \
  && rm -rf /var/lib/apt/lists/*

ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA256 43ad621c9b014cde8db17393dc108378d37bc853aa351a6c74bf6432c1bbd182

RUN curl -fsSL --insecure "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
  && echo "$GOLANG_DOWNLOAD_SHA256 golang.tar.gz" | sha256sum -c - \
  && tar -C /usr/local -xzf golang.tar.gz \
  && rm golang.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH

# Fetch the latest version of the package dependencies
RUN go get -u \
  golang.org/x/net/context \
  github.com/daaku/go.httpgzip \
  github.com/rs/cors \
  github.com/tj/go-debug \
  gopkg.in/h2non/filetype.v0 \
  gopkg.in/throttled/throttled.v2 \
  gopkg.in/throttled/throttled.v2/store/memstore

# Fetch and patch bimg
COPY bimg-resize.patch /tmp
RUN git clone https://github.com/h2non/bimg.git /go/src/gopkg.in/h2non/bimg.v1 && \
    patch /go/src/gopkg.in/h2non/bimg.v1/resize.go /tmp/bimg-resize.patch

# Build imaginary
COPY *.go /go/src/imaginary/
RUN cd /go/src/imaginary && \
    go build -o /go/bin/imaginary && \
    cd /go && \
    rm -rf /go/src

# Run the outyet command by default when the container starts.
ENTRYPOINT ["/go/bin/imaginary"]

# Expose the server TCP port
EXPOSE 9000
