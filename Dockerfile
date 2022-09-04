FROM golang:1.19.0-bullseye as go-builder

ARG IMAGESYNC_VERSION
ENV IMAGESYNC_VERSION ${IMAGESYNC_VERSION}

ENV PACKAGE github.com/MQasimSarfraz/image-sync
ENV CGO_ENABLED 1

WORKDIR $GOPATH/src/$PACKAGE

# create directories for binary and install dependencies
RUN mkdir -p /out && \
    apt -qq update && \
    apt install -y git libgpgme-dev libassuan-dev libbtrfs-dev libdevmapper-dev pkg-config

# copy sources, test and build the application
COPY . ./
RUN go vet ./...
RUN go test --parallel=1 ./...
RUN go build -v \
    -ldflags="-s -w \
      -X 'github.com/MQasimSarfraz/imagesync.Version=${IMAGESYNC_VERSION}'"\
    -o /out/imagesync ./cmd/imagesync


# build the final container image
FROM bitnami/minideb:bullseye

RUN apt -qq update && \
    apt install -y libgpgme-dev --no-install-recommends && \
    rm -rf /var/lib/apt/lists/*

COPY --from=go-builder /out/imagesync /

ENTRYPOINT ["/imagesync"]

CMD ["-h"]
