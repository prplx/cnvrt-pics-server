ARG GOLANG_VERSION=1.22.5

FROM golang:${GOLANG_VERSION}-bookworm as builder

ARG VIPS_VERSION=8.15.2
ARG CGIF_VERSION=0.4.1
ARG LIBSPNG_VERSION=0.7.4
ARG TARGETARCH=arm

ENV PKG_CONFIG_PATH=/usr/local/lib/pkgconfig

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
  libgsf-1-dev libfftw3-dev liborc-0.4-dev librsvg2-dev libimagequant-dev libheif-dev libgirepository1.0-dev && \
  pip3 install meson --break-system-packages && \
  cd /tmp && \
    curl -fsSLO https://github.com/dloebl/cgif/archive/refs/tags/v${CGIF_VERSION}.tar.gz && \
    tar xf v${CGIF_VERSION}.tar.gz && \
    cd cgif-${CGIF_VERSION} && \
    meson build --prefix=/usr/local --libdir=/usr/local/lib --buildtype=release && \
    cd build && \
    ninja && \
    ninja install && \
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
    curl -fsSLO https://github.com/libvips/libvips/releases/download/v${VIPS_VERSION}/vips-${VIPS_VERSION}.tar.xz && \
    tar xf vips-${VIPS_VERSION}.tar.xz && \
    cd vips-${VIPS_VERSION} && \
    meson setup build --prefix /usr/local && \
    cd build && \
    meson compile && \
    meson test && \
    meson install && \
    ldconfig && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-${TARGETARCH}64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate && \
    rm -rf /usr/local/lib/python* && \
    rm -rf /usr/local/lib/libvips-cpp.* && \
    rm -rf /usr/local/lib/*.a && \
    rm -rf /usr/local/lib/*.la

WORKDIR ${GOPATH}/src/github.com/prplx/cnvrt

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# RUN touch .envrc && make test
RUN go build -o ${GOPATH}/bin/cnvrt ./cmd/api/main.go

FROM debian:bookworm-slim as prod

COPY --from=builder /usr/local/lib /usr/local/lib
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate
COPY --from=builder /go/bin/cnvrt /usr/local/bin/cnvrt
COPY config.yaml /app/
COPY migrations /app/migrations

# Install runtime dependencies
RUN DEBIAN_FRONTEND=noninteractive \
apt-get update && \
apt-get install --no-install-recommends -y \
curl procps libglib2.0-0 libjpeg62-turbo libpng16-16 libopenexr-3-1-30 \
libwebp7 libwebpmux3 libwebpdemux2 libtiff6 libexif12 libxml2 libpoppler-glib8 \
libpango1.0-0 libmatio11 libopenslide0 libopenjp2-7 libjemalloc2 \
libgsf-1-114 libfftw3-bin liborc-0.4-0 librsvg2-2 libcfitsio10 libimagequant0 dav1d libheif1 && \
ln -s /usr/lib/$(uname -m)-linux-gnu/libjemalloc.so.2 /usr/local/lib/libjemalloc.so && \
apt-get autoremove -y && \
apt-get autoclean && \
apt-get clean && \
rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

ENV VIPS_WARNING=0
ENV MALLOC_ARENA_MAX=2
ENV LD_PRELOAD=/usr/local/lib/libjemalloc.so

RUN chown -R nobody:nogroup /app
RUN chmod 755 /app

USER nobody

EXPOSE ${PORT}

CMD ["/app/start.sh"]

FROM builder AS dev

RUN go install github.com/air-verse/air@latest

COPY --from=builder /go/pkg /go/pkg

EXPOSE ${PORT}

CMD ["/app/start.sh"]
