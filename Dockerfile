ARG GOLANG_VERSION=1.22.5
FROM golang:${GOLANG_VERSION}-bookworm as builder

ARG VIPS_VERSION=8.15.2
ARG CGIF_VERSION=0.4.1
ARG LIBSPNG_VERSION=0.7.4
ARG TARGETARCH

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
    meson setup build --prefix /usr/local && \
    cd build && \
    meson compile && \
    meson test && \
    meson install && \
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

# RUN touch .envrc && make test
RUN go build -o ${GOPATH}/bin/cnvrt ./cmd/api/main.go

FROM debian:bookworm-slim as base

COPY --from=builder /usr/local/lib /usr/local/lib
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

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
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-arm64.tar.gz | tar xvz && \
mv migrate /usr/local/bin/migrate && \
rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY --from=builder /go/bin/cnvrt /usr/local/bin/cnvrt
COPY config.yaml /app/
COPY migrations /app/migrations

ENV VIPS_WARNING=0
ENV MALLOC_ARENA_MAX=2
ENV LD_PRELOAD=/usr/local/lib/libjemalloc.so

FROM base AS prod

RUN chown -R nobody:nogroup /app
RUN chmod 755 /app

# use unprivileged user
USER nobody

EXPOSE ${PORT}

CMD ["/app/start.sh"]

FROM builder AS dev

ARG GO_VERSION=1.22.5

WORKDIR /app

# Download and install Go
RUN apt-get update && apt-get install make && \
curl -L https://golang.org/dl/go${GO_VERSION}.linux-arm64.tar.gz -o go${GO_VERSION}.linux-arm64.tar.gz \
&& tar -C /usr/local -xzf go${GO_VERSION}.linux-arm64.tar.gz \
&& rm go${GO_VERSION}.linux-arm64.tar.gz

# Set Go environment variables
ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV PATH $GOPATH/bin:$GOROOT/bin:$PATH

# Create Go workspace directory
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

RUN go install github.com/air-verse/air@latest

COPY --from=builder /go/pkg /go/pkg
COPY --from=base /usr/local/bin/migrate /usr/local/bin/migrate

EXPOSE 3002

CMD ["/app/start.sh"]
