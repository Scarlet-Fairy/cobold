#syntax=docker/dockerfile:1.2

FROM golang:1.15.8-alpine AS base

MAINTAINER Michele Della Mea <michele.dellamea.arcanediver@gmail.com>

# Create appuser.
ARG USER=appuser
ARG UID=10001

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR /src

ENV CGO_ENABLED=0 \
    GO111MODULE=on

COPY go.* .
RUN go mod download
COPY . .

# ---------------------- #

FROM base AS build

ARG TARGETOS=linux
ARG TARGETARCH=amd64

RUN -mount=type=cache,target=/root/.cache/go-build \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/cobold ./cmd/cobold/main.go

# ---------------------- #

FROM scratch


COPY --from=base /etc/passwd /etc/passwd
COPY --from=base /etc/group /etc/group

COPY --from=build /out/cobold .

USER appuser:appuser

ENTRYPOINT ["/cobold"]