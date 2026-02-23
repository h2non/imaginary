# syntax=docker/dockerfile:1.4

ARG GOLANG_VERSION=1.25
FROM golang:${GOLANG_VERSION}-bullseye as builder


ARG IMAGINARY_VERSION=dev
ARG LIBVIPS_VERSION=8.14.1
ARG GOLANGCILINT_VERSION=1.51.1

ENV LIBSPNG_VERSION="0.7.3"
ENV LIBSPNG_URL="https://github.com/randy408/libspng/archive/refs/tags/v${LIBSPNG_VERSION}.tar.gz"

#ENV PDFIUM_VERSION="5579"
#ENV PDFIUM_URL="https://github.com/bblanchon/pdfium-binaries/releases/download/chromium/${PDFIUM_VERSION}/pdfium-linux-arm64.tgz"
#
## Installs libvips + required libraries
#COPY <<EOF /usr/lib/pkgconfig/pdfium.pc
#prefix=/usr
#exec_prefix=\${prefix}
#includedir=\${prefix}/include
#libdir=\${exec_prefix}/lib/aarch64-linux-gnu
#
#Name: pdfium
#Description: pdfium
#Version: ${PDFIUM_VERSION}
#Cflags: -I\${includedir}
#Libs: -L\${libdir} -lpdfium
#EOF

RUN DEBIAN_FRONTEND=noninteractive \
  apt-get update && \
  apt-get install --no-install-recommends -y \
  ca-certificates \
  automake build-essential curl meson scons file \
  gobject-introspection gtk-doc-tools libglib2.0-dev libjpeg62-turbo-dev libpng-dev \
  libwebp-dev libtiff5-dev libgif-dev libexif-dev libxml2-dev libpoppler-glib-dev \
  swig libmagickwand-dev libpango1.0-dev libmatio-dev libopenslide-dev libcfitsio-dev \
  libgsf-1-dev fftw3-dev liborc-0.4-dev librsvg2-dev libimagequant-dev libheif-dev \
  libexpat1-dev libgirepository1.0-dev libglib2.0-dev liblcms2-dev libnifti2-dev libniftiio-dev libopenexr-dev \
  libopenjp2-7-dev patchelf pkg-config && \
  cd /tmp && \
  curl -L "${LIBSPNG_URL}" | tar -xzf- && \
  cd /tmp/libspng-${LIBSPNG_VERSION} && \
  meson build --buildtype=release && \
  cd /tmp/libspng-${LIBSPNG_VERSION}/build && \
  ninja && ninja install && \
#  cd /usr && \
#  curl -L "${PDFIUM_URL}" | tar -xzf- && \
  cd /tmp && \
  curl -fsSLO https://github.com/libvips/libvips/releases/download/v${LIBVIPS_VERSION}/vips-${LIBVIPS_VERSION}.tar.xz && \
  tar xJf vips-${LIBVIPS_VERSION}.tar.xz && \
  cd /tmp/vips-${LIBVIPS_VERSION} && \
  CFLAGS="-g -O3" CXXFLAGS="-D_GLIBCXX_USE_CXX11_ABI=0 -g -O3" meson setup build-dir --buildtype=release && \
  cd /tmp/vips-${LIBVIPS_VERSION}/build-dir && \
  CFLAGS="-g -O3" CXXFLAGS="-D_GLIBCXX_USE_CXX11_ABI=0 -g -O3" meson compile && \
  CFLAGS="-g -O3" CXXFLAGS="-D_GLIBCXX_USE_CXX11_ABI=0 -g -O3" meson install && \
  ldconfig

# Installing golangci-lint
WORKDIR /tmp
RUN curl -fsSL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "${GOPATH}/bin" v${GOLANGCILINT_VERSION}

WORKDIR ${GOPATH}/src/github.com/h2non/imaginary

# Cache go modules
ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .

RUN go mod download

# Copy imaginary sources
COPY . .

# Run quality control
RUN go test ./... -test.v -race -test.coverprofile=atomic .
RUN golangci-lint run .

# Compile imaginary
RUN go build -a \
    -o ${GOPATH}/bin/imaginary \
    -ldflags="-s -w -h -X main.Version=${IMAGINARY_VERSION}" \
    github.com/h2non/imaginary

FROM debian:bullseye-slim

ARG IMAGINARY_VERSION

LABEL maintainer="tomas@aparicio.me" \
      org.label-schema.description="Fast, simple, scalable HTTP microservice for high-level image processing with first-class Docker support" \
      org.label-schema.schema-version="1.0" \
      org.label-schema.url="https://github.com/h2non/imaginary" \
      org.label-schema.vcs-url="https://github.com/h2non/imaginary" \
      org.label-schema.version="${IMAGINARY_VERSION}"

COPY --from=builder /usr/local/lib /usr/local/lib
COPY --from=builder /go/bin/imaginary /usr/local/bin/imaginary
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

# Install runtime dependencies
RUN DEBIAN_FRONTEND=noninteractive \
  apt-get update && \
  apt-get install --no-install-recommends -y \
  procps libglib2.0-0 libjpeg62-turbo libpng16-16 libopenexr25 \
  libwebp6 libwebpmux3 libwebpdemux2 libtiff5 libgif7 libexif12 libxml2 libpoppler-glib8 \
  libmagickwand-6.q16-6 libpango1.0-0 libmatio11 libopenslide0 libjemalloc2 \
  libgsf-1-114 fftw3 liborc-0.4-0 librsvg2-2 libcfitsio9 libimagequant0 libheif1 \
  liblcms2-2 libopenexr25 libopenjp2-7 libnifti2-2 libniftiio2 && \
  ln -s /usr/lib/$(uname -m)-linux-gnu/libjemalloc.so.2 /usr/local/lib/libjemalloc.so && \
  apt-get autoremove -y && \
  apt-get autoclean && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
ENV LD_PRELOAD=/usr/local/lib/libjemalloc.so

# Server port to listen
ENV PORT=9000

# Drop privileges for non-UID mapped environments
USER nobody

VOLUME /config

# Run the entrypoint command by default when the container starts.
ENTRYPOINT ["/usr/local/bin/imaginary"]

# Expose the server TCP port
EXPOSE ${PORT}
