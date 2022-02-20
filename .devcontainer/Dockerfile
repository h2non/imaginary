# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.217.4/containers/go/.devcontainer/base.Dockerfile

# [Choice] Go version (use -bullseye variants on local arm64/Apple Silicon): 1, 1.16, 1.17, 1-bullseye, 1.16-bullseye, 1.17-bullseye, 1-buster, 1.16-buster, 1.17-buster
ARG VARIANT="1.17-bullseye"
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT}

# Versions of libvips and golanci-lint
ARG LIBVIPS_VERSION=8.12.2
ARG GOLANGCILINT_VERSION=1.29.0

# Install additional OS packages
RUN DEBIAN_FRONTEND=noninteractive \
  apt-get update && \
  apt-get install --no-install-recommends -y \
  ca-certificates \
  automake build-essential curl \
  procps    libopenexr25 libmagickwand-6.q16-6 libpango1.0-0 libmatio11 \
  libopenslide0 libjemalloc2 gobject-introspection gtk-doc-tools \
  libglib2.0-0 libglib2.0-dev libjpeg62-turbo libjpeg62-turbo-dev \
  libpng16-16 libpng-dev libwebp6 libwebpmux3 libwebpdemux2 libwebp-dev \
  libtiff5 libtiff5-dev libgif7 libgif-dev libexif12 libexif-dev \
  libxml2 libxml2-dev libpoppler-glib8 libpoppler-glib-dev \
  swig libmagickwand-dev libpango1.0-dev libmatio-dev libopenslide-dev \
  libcfitsio9 libcfitsio-dev libgsf-1-114 libgsf-1-dev fftw3 fftw3-dev \
  liborc-0.4-0 liborc-0.4-dev librsvg2-2 librsvg2-dev libimagequant0 \
  libimagequant-dev libheif1 libheif-dev && \
  cd /tmp && \
  curl -fsSLO https://github.com/libvips/libvips/releases/download/v${LIBVIPS_VERSION}/vips-${LIBVIPS_VERSION}.tar.gz && \
  tar zvxf vips-${LIBVIPS_VERSION}.tar.gz && \
  cd /tmp/vips-${LIBVIPS_VERSION} && \
	CFLAGS="-g -O3" CXXFLAGS="-D_GLIBCXX_USE_CXX11_ABI=0 -g -O3" \
    ./configure \
    --disable-debug \
    --disable-dependency-tracking \
    --disable-introspection \
    --disable-static \
    --enable-gtk-doc-html=no \
    --enable-gtk-doc=no \
    --enable-pyvips8=no && \
  make && \
  make install && \
  ldconfig

# Installing golangci-lint
RUN curl -fsSL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "${GOPATH}/bin" v${GOLANGCILINT_VERSION}

# [Optional] Uncomment the next lines to use go get to install anything else you need
# USER vscode
# RUN go get -x <your-dependency-or-tool>
