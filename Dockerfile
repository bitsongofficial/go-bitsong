# To build the BitSong image, just run:
# > docker build -t bitsongofficial/go-bitsong .
#
# Simple usage with a mounted data directory:
# > docker run -it -p 26657:26657 -p 26656:26656 -v ~/.bitsongd:/root/.bitsongd bitsongofficial/go-bitsong bitsongd init
# > docker run -it -p 26657:26657 -p 26656:26656 -v ~/.bitsongd:/root/.bitsongd bitsongofficial/go-bitsong bitsongd start

FROM golang:1.17.5-alpine AS build-env

# Set up dependencies
RUN apk add --no-cache curl make git libc-dev bash gcc linux-headers eudev-dev python3

# Set working directory for the build
WORKDIR /go/src/github.com/bitsongofficial/go-bitsong

# Add source files
COPY . .

# Build BitSong
RUN make build-linux

# Final image
FROM alpine:edge

# Install ca-certificates
RUN apk add --update ca-certificates
WORKDIR /root

# Install bash
RUN apk add --no-cache bash

# Copy over binaries from the build-env
COPY --from=build-env /go/src/github.com/bitsongofficial/go-bitsong/build/bitsongd /usr/bin/bitsongd

EXPOSE 26656 26657 1317 9090

CMD ["bitsongd"]