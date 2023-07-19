ARG VARIANT="1.20-bullseye"
FROM mcr.microsoft.com/devcontainers/go:0-${VARIANT}

# Versions of libvips
ARG LIBVIPS_VERSION=8.14.2

# Install additional OS packages
RUN DEBIAN_FRONTEND=noninteractive \
  apt-get update && \
  apt-get install --no-install-recommends -y \
  ca-certificates \
  automake build-essential curl xz-utils pkg-config libexpat1-dev \
  python3 python3-pip ninja-build libgirepository1.0-dev \
  procps libopenexr25 libmagickwand-6.q16-6 libpango1.0-0 libmatio11 \
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
  curl -fsSLO https://github.com/libvips/libvips/releases/download/v${LIBVIPS_VERSION}/vips-${LIBVIPS_VERSION}.tar.xz && \
  tar vxf vips-${LIBVIPS_VERSION}.tar.xz && \
  cd /tmp/vips-${LIBVIPS_VERSION} && \
  pip3 install meson && \
  meson setup build --prefix /usr/local && \
  cd build && \
  meson compile && \
  meson test && \
  meson install && \
  ldconfig

# [Optional] Uncomment the next lines to use go get to install anything else you need
# USER vscode
# RUN go get -x <your-dependency-or-tool>
