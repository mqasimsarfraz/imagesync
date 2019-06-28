FROM golang:1.12.6-stretch as go-builder

ENV PACKAGE github.com/MQasimSarfraz/image-sync
ENV CGO_ENABLED 1
ENV GO111MODULE=on

WORKDIR $GOPATH/src/$PACKAGE

# create directories for binary and install dependencies
RUN mkdir -p /out && \
    apt -qq update && \
    apt install -y git libgpgme-dev libostree-dev

# copy sources, test and build the application
COPY . ./
RUN go vet ./...
RUN go test --parallel=1 ./...
RUN go build -v -ldflags="-s -w" -o /out/imagesync ./cmd/imagesync


# build the final container image
FROM bitnami/minideb:stretch

RUN apt -qq update && \
    apt install -y git libgpgme-dev && \
    rm -rf /var/lib/apt/lists/*

EXPOSE 3080

COPY --from=go-builder /out/imagesync /

ENTRYPOINT ["/imagesync"]

CMD ["-h"]