# syntax=docker/dockerfile:1

ARG BASE_IMG_TAG=nonroot

# --------------------------------------------------------
# Build 
# --------------------------------------------------------

FROM golang:1.22-alpine as build

RUN set -eux; apk add --no-cache ca-certificates build-base;
RUN apk add git
# Needed by github.com/zondax/hid
RUN apk add linux-headers

WORKDIR /src/app/
COPY go.mod go.sum* ./
RUN go mod download
COPY . .

ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev
RUN apk add --no-cache $PACKAGES

RUN LEDGER_ENABLED=false CGO_ENABLED=0 make install

# --------------------------------------------------------
# Runner
# --------------------------------------------------------

FROM gcr.io/distroless/base-debian11:${BASE_IMG_TAG}
#FROM ubuntu:20.04

ENV HOME /go-bitsong
WORKDIR $HOME

COPY --from=build /go/bin/bitsongd /bin/bitsongd

EXPOSE 26656
EXPOSE 26657
EXPOSE 9090
EXPOSE 1317

#ENTRYPOINT ["/bin/bash"]
ENTRYPOINT ["bitsongd"]
CMD [ "start" ]