# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM marcbachmann/libvips

# Go version to use
ENV GOLANG_VERSION 1.4.2

# Server port to listen
ENV PORT 9000

# Install dependencies
RUN apt-get update -y && \
  apt-get install -y curl git

# gcc for cgo
RUN apt-get install -y \
    gcc libc6-dev make \
    --no-install-recommends \
  && rm -rf /var/lib/apt/lists/*

RUN curl -sSL https://golang.org/dl/go$GOLANG_VERSION.src.tar.gz \
    | tar -v -C /usr/src -xz

RUN cd /usr/src/go/src && ./make.bash --no-clean 2>&1

ENV PATH /usr/src/go/bin:$PATH

RUN mkdir -p /go/src /go/bin && chmod -R 777 /go
ENV GOPATH /go
ENV PATH /go/bin:$PATH
WORKDIR /go

# Fetch the latest version of the package
RUN go get -u github.com/h2non/imaginary

# Run the outyet command by default when the container starts.
ENTRYPOINT ["/go/bin/imaginary"]
CMD ["-help"]

# Expose the server TCP port
EXPOSE 9000
