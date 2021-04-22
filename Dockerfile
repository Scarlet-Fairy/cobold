FROM golang:1.15.8-alpine AS base

LABEL org.opencontainers.image.source=https://github.com/Scarlet-Fairy/cobold

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

# ENV CGO_ENABLED=0 \
#     GO111MODULE=on

COPY go.* ./
RUN go mod download
COPY . .

# ---------------------- #

FROM base AS build

ENV TARGETOS=linux
ENV TARGETARCH=amd64

RUN GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -mod vendor \
    -o /out/cobold \
    ./cmd/cobold/main.go

# ---------------------- #

FROM alpine:latest

RUN apk update
RUN apk add ca-certificates git
RUN rm -rf /var/cache/apk/*

#COPY --from=base /etc/passwd /etc/passwd
#COPY --from=base /etc/group /etc/group

COPY --from=build /out/cobold .

#USER appuser:appuser

ENTRYPOINT ["/cobold"]