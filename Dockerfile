# docker build . -t bitsongofficial/go-bitsong:latest
# docker run --rm -it bitsongofficial/go-bitsong:latest /bin/sh
FROM golang:1.23-alpine AS go-builder

# this comes from standard alpine nightly file
#  https://github.com/rust-lang/docker-rust-nightly/blob/master/alpine3.12/Dockerfile
# with some changes to support our toolchain, etc
SHELL ["/bin/sh", "-ecuxo", "pipefail"]
# we probably want to default to latest and error
# since this is predominantly for dev use
# hadolint ignore=DL3018
RUN apk add --no-cache ca-certificates build-base git
# NOTE: add these to run with LEDGER_ENABLED=true
# RUN apk add libusb-dev linux-headers

WORKDIR /code

# Download dependencies and CosmWasm libwasmvm if found.
ADD go.mod go.sum ./

# Cosmwasm - Download correct libwasmvm version
RUN ARCH=$(uname -m) && WASMVM_VERSION=$(go list -m github.com/CosmWasm/wasmvm/v2 | sed 's/.* //') && \
    wget https://github.com/CosmWasm/wasmvm/releases/download/$WASMVM_VERSION/libwasmvm_muslc.$ARCH.a \
    -O /lib/libwasmvm_muslc.$ARCH.a && \
    # verify checksum
    wget https://github.com/CosmWasm/wasmvm/releases/download/$WASMVM_VERSION/checksums.txt -O /tmp/checksums.txt && \
    sha256sum /lib/libwasmvm_muslc.$ARCH.a | grep $(cat /tmp/checksums.txt | grep libwasmvm_muslc.$ARCH | cut -d ' ' -f 1) 

# Copy over code
COPY . /code/

# force it to use static lib (from above) not standard libgo_cosmwasm.so file
# then log output of file /code/build/bitsongd
# then ensure static linking
RUN LEDGER_ENABLED=false BUILD_TAGS=muslc LINK_STATICALLY=true make build \
  && file /code/build/bitsongd \
  && echo "Ensuring binary is statically linked ..." \
  && (file /code/build/bitsongd | grep "statically linked")

# --------------------------------------------------------
FROM alpine:3.16

COPY --from=go-builder /code/build/bitsongd /usr/bin/bitsongd

ENV HOME=/bitsongd
WORKDIR $HOME

# rest server, tendermint p2p, tendermint rpc
EXPOSE 1317 26656 26657

CMD ["/usr/bin/bitsongd"]

## To use as CLI:
# alias bsd="docker run --rm -it -v ~/.bitsongd:/root/.bitsongd bitsong:v0.12323 bitsongd"