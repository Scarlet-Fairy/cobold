#syntax=docker/dockerfile:1.2

FROM golang:1.15.8-alpine AS base

WORKDIR /src

ENV CGO_ENABLED=0

COPY go.* .
RUN go mod download
COPY . .

# ---------------------- #

FROM base AS build

ARG TARGETOS=linux
ARG TARGETARCH=amd64

RUN -mount=type=cache,target=/root/.cache/go-build \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/cobold .

# ---------------------- #

FROM scratch
COPY --from=build /out/cobold .
ENTRYPOINT ["/cobold"]