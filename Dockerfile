ARG GOLANG_VERSION=1.20.8
FROM golang:${GOLANG_VERSION}-bullseye as builder

ARG VIPS_VERSION=8.14.5
ARG CGIF_VERSION=0.3.0
ARG LIBSPNG_VERSION=0.7.3
ARG TARGETARCH
ARG db_dsn

ENV PKG_CONFIG_PATH=/usr/local/lib/pkgconfig
ENV DB_DSN=$db_dsn

# libaom3 is in Debian bullseye-backports
RUN echo 'deb http://deb.debian.org/debian bullseye-backports main' > /etc/apt/sources.list.d/backports.list

# Installs libvips + required libraries
RUN DEBIAN_FRONTEND=noninteractive \ 
  apt-get update && \
  apt-get install --no-install-recommends -y \
  ca-certificates \
  automake build-essential curl \
  python3-pip ninja-build pkg-config \
  gobject-introspection gtk-doc-tools libglib2.0-dev libjpeg62-turbo-dev libpng-dev \
  libwebp-dev libtiff5-dev libexif-dev libxml2-dev libpoppler-glib-dev \
  swig libpango1.0-dev libmatio-dev libopenslide-dev libcfitsio-dev libopenjp2-7-dev liblcms2-dev \
  libgsf-1-dev fftw3-dev liborc-0.4-dev librsvg2-dev libimagequant-dev libaom-dev/bullseye-backports libheif-dev && \
  pip3 install meson && \
  cd /tmp && \
    curl -fsSLO https://github.com/dloebl/cgif/archive/refs/tags/V${CGIF_VERSION}.tar.gz && \
    tar xf V${CGIF_VERSION}.tar.gz && \
    cd cgif-${CGIF_VERSION} && \
    meson build --prefix=/usr/local --libdir=/usr/local/lib --buildtype=release && \
    cd build && \
    ninja && \
    ninja install && \
  cd /tmp && \
    curl -fsSLO https://github.com/randy408/libspng/archive/refs/tags/v${LIBSPNG_VERSION}.tar.gz && \
    tar xf v${LIBSPNG_VERSION}.tar.gz && \
    cd libspng-${LIBSPNG_VERSION} && \
    meson setup _build \
      --buildtype=release \
      --strip \
      --prefix=/usr/local \
      --libdir=lib && \
    ninja -C _build && \
    ninja -C _build install && \
  cd /tmp && \
    curl -fsSLO https://github.com/libvips/libvips/releases/download/v${VIPS_VERSION}/vips-${VIPS_VERSION}.tar.xz && \
    tar xf vips-${VIPS_VERSION}.tar.xz && \
    cd vips-${VIPS_VERSION} && \
    meson setup _build \
    --buildtype=release \
    --strip \
    --prefix=/usr/local \
    --libdir=lib \
    -Dgtk_doc=false \
    -Dmagick=disabled \
    -Dintrospection=false && \
    ninja -C _build && \
    ninja -C _build install && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate $GOPATH/bin/migrate && \
  ldconfig && \
  rm -rf /usr/local/lib/python* && \
  rm -rf /usr/local/lib/libvips-cpp.* && \
  rm -rf /usr/local/lib/*.a && \
  rm -rf /usr/local/lib/*.la

WORKDIR ${GOPATH}/src/github.com/prplx/cnvrt

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN if [ "$TARGETARCH" = "amd64" ]; then touch .envrc && make test; fi
RUN go build -o ${GOPATH}/bin/cnvrt ./cmd/api/main.go
RUN make db/migrate_up

FROM debian:bullseye-slim

COPY --from=builder /usr/local/lib /usr/local/lib
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /etc/apt/sources.list.d/backports.list /etc/apt/sources.list.d/backports.list 

# Install runtime dependencies
RUN DEBIAN_FRONTEND=noninteractive \
  apt-get update && \
  apt-get install --no-install-recommends -y \
  procps libglib2.0-0 libjpeg62-turbo libpng16-16 libopenexr25 \
  libwebp6 libwebpmux3 libwebpdemux2 libtiff5 libexif12 libxml2 libpoppler-glib8 \
  libpango1.0-0 libmatio11 libopenslide0 libopenjp2-7 libjemalloc2 \
  libgsf-1-114 fftw3 liborc-0.4-0 librsvg2-2 libcfitsio9 libimagequant0 libaom3 libheif1 && \
  ln -s /usr/lib/$(uname -m)-linux-gnu/libjemalloc.so.2 /usr/local/lib/libjemalloc.so && \
  apt-get autoremove -y && \
  apt-get autoclean && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY --from=builder /go/bin/cnvrt /usr/local/bin/cnvrt
COPY config.yaml /app/config.yaml

ENV VIPS_WARNING=0
ENV MALLOC_ARENA_MAX=2
ENV LD_PRELOAD=/usr/local/lib/libjemalloc.so

RUN chown -R nobody:nogroup /app
RUN chmod 755 /app
# use unprivileged user
USER nobody

CMD /usr/local/bin/cnvrt \
  -env=${ENV} \
  -port=${PORT} \
  -upload-dir=${UPLOAD_DIR} \
  -db-dsn=${DB_DSN} \
  -metrics-user=${METRICS_USER} \
  -metrics-password=${METRICS_PASSWORD} \
  -firebase-project-id=${FIREBASE_PROJECT_ID} \
  -allow-origins=${ALLOW_ORIGINS}

EXPOSE ${PORT}
