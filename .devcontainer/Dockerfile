FROM mcr.microsoft.com/devcontainers/go:1-1.21-bullseye

# Install additional OS packages
RUN DEBIAN_FRONTEND=noninteractive \
  apt-get update && \
  apt-get install --no-install-recommends -y \
  ca-certificates \
  curl \
  git \
  build-essential \
  libvips-dev \
  libmagickwand-dev \
  libpng-dev \
  libjpeg-dev \
  libwebp-dev \
  libtiff-dev \
  libgif-dev \
  libexif-dev \
  libxml2-dev \
  libpoppler-glib-dev \
  libopenslide-dev \
  libcfitsio-dev \
  libgsf-1-dev \
  fftw3-dev \
  liborc-0.4-dev \
  librsvg2-dev \
  libimagequant-dev \
  libheif-dev \
  libopenexr-dev \
  libmatio-dev && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/*

# Install additional Go tools
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest && \
    go install golang.org/x/tools/cmd/goimports@latest && \
    go install github.com/fatih/gomodifytags@latest

# Set up workspace
WORKDIR /workspace
