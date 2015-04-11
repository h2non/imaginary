# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM marcbachmann/libvips

# gcc for cgo
RUN apt-get update && apt-get install -y \
    gcc libc6-dev make \
    --no-install-recommends \
  && rm -rf /var/lib/apt/lists/*

ENV GOLANG_VERSION 1.4.2

RUN curl -sSL https://golang.org/dl/go$GOLANG_VERSION.src.tar.gz \
    | tar -v -C /usr/src -xz

RUN cd /usr/src/go/src && ./make.bash --no-clean 2>&1

ENV PATH /usr/src/go/bin:$PATH

RUN mkdir -p /go/src /go/bin && chmod -R 777 /go
ENV GOPATH /go
ENV PATH /go/bin:$PATH
WORKDIR /go

COPY go-wrapper /usr/local/bin/

# Server port to listen
ENV PORT 80
# Enable debug mode?
#ENV DEBUG *

# Fetch the latest version of the package
RUN go get gopkg.in/h2non/imaginary.v0

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/imaginary

# Expose the server TCP port
EXPOSE 80
