# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.217.4/containers/go/.devcontainer/base.Dockerfile

# [Choice] Go version (use -bullseye variants on local arm64/Apple Silicon): 1, 1.16, 1.17, 1-bullseye, 1.16-bullseye, 1.17-bullseye, 1-buster, 1.16-buster, 1.17-buster
ARG VARIANT="1.19-bullseye"
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT}

# Version of golangci-lint
ARG GOLANGCILINT_VERSION=1.29.0

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
  libheif-dev

# Installing golangci-lint
RUN curl -fsSL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "${GOPATH}/bin" v${GOLANGCILINT_VERSION}

# [Optional] Uncomment the next lines to use go get to install anything else you need
# USER vscode
# RUN go get -x <your-dependency-or-tool>
