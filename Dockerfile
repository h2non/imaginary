# Build Stage
ARG GOLANG_VERSION=1.17
FROM --platform=$TARGETARCH golang:${GOLANG_VERSION}-bullseye as builder

# Set GO environment variables for cross-compilation
ENV GOARCH=amd64 

ARG IMAGINARY_VERSION=1.2.5
ARG LIBVIPS_VERSION=8.12.2
ARG GOLANGCILINT_VERSION=1.29.0

# Install dependencies
RUN apt-get update && apt-get install --no-install-recommends -y \
  ca-certificates automake build-essential curl \
  gobject-introspection gtk-doc-tools libglib2.0-dev libjpeg62-turbo-dev libpng-dev \
  libwebp-dev libtiff5-dev libgif-dev libexif-dev libxml2-dev libpoppler-glib-dev \
  swig libmagickwand-dev libpango1.0-dev libmatio-dev libopenslide-dev libcfitsio-dev \
  libgsf-1-dev fftw3-dev liborc-0.4-dev librsvg2-dev libimagequant-dev libheif-dev

# Install libvips (compiled from source)
WORKDIR /tmp
RUN curl -fsSLO https://github.com/libvips/libvips/releases/download/v${LIBVIPS_VERSION}/vips-${LIBVIPS_VERSION}.tar.gz && \
  tar xzf vips-${LIBVIPS_VERSION}.tar.gz && cd vips-${LIBVIPS_VERSION} && \
  CFLAGS="-g -O1" CXXFLAGS="-D_GLIBCXX_USE_CXX11_ABI=0 -g -O3" \
  ./configure \
  --disable-debug \
  --disable-dependency-tracking \
  --disable-introspection \
  --disable-static \
  --enable-gtk-doc-html=no \
  --enable-gtk-doc=no \
  --enable-pyvips8=no && \
  make -j$(nproc) && make install && ldconfig

# Install golangci-lint
WORKDIR /tmp
RUN curl -fsSL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "${GOPATH}/bin" v${GOLANGCILINT_VERSION}

# Set up imaginary
WORKDIR ${GOPATH}/src/github.com/h2non/imaginary
COPY go.mod go.sum ./
RUN go mod download && go mod tidy

# Fix missing dependencies
RUN go get github.com/throttled/throttled/v2
RUN go get github.com/throttled/throttled/v2/store/memstore@v2.13.0

# Copy source files
COPY . .

# Run tests & lint
RUN go test ./... -test.v -race -test.coverprofile=atomic .
RUN golangci-lint run .

# Compile imaginary binary
RUN go build -a -o ${GOPATH}/bin/imaginary \
  -ldflags="-s -w -X main.Version=${IMAGINARY_VERSION}" \
  github.com/h2non/imaginary

# Final Stage (Slim Runtime)
FROM --platform=linux/amd64 debian:bullseye-slim

ARG IMAGINARY_VERSION

LABEL maintainer="tomas@aparicio.me" \
  org.label-schema.description="Fast, simple, scalable HTTP microservice for high-level image processing with first-class Docker support" \
  org.label-schema.schema-version="1.0" \
  org.label-schema.url="https://github.com/h2non/imaginary" \
  org.label-schema.vcs-url="https://github.com/h2non/imaginary" \
  org.label-schema.version="${IMAGINARY_VERSION}"

# Copy compiled binary & dependencies from builder
COPY --from=builder /usr/local/lib /usr/local/lib
COPY --from=builder /go/bin/imaginary /usr/local/bin/imaginary
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

# Install minimal runtime dependencies
RUN apt-get update && apt-get install --no-install-recommends -y \
  procps libglib2.0-0 libjpeg62-turbo libpng16-16 libopenexr25 \
  libwebp6 libwebpmux3 libwebpdemux2 libtiff5 libgif7 libexif12 libxml2 libpoppler-glib8 \
  libmagickwand-6.q16-6 libpango1.0-0 libmatio11 libopenslide0 libjemalloc2 \
  libgsf-1-114 fftw3 liborc-0.4-0 librsvg2-2 libcfitsio9 libimagequant0 libheif1 && \
  ln -s /usr/lib/$(uname -m)-linux-gnu/libjemalloc.so.2 /usr/local/lib/libjemalloc.so && \
  apt-get autoremove -y && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Configure memory optimization
ENV LD_PRELOAD=/usr/local/lib/libjemalloc.so

# Expose server port
ENV PORT 9000
EXPOSE ${PORT}

# Run as non-root user
USER nobody

# Start the service
ENTRYPOINT ["/usr/local/bin/imaginary"]
