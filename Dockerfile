# base
FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.17 as base

WORKDIR /go/src/github.com/shteou/gwrp

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ cmd
COPY pkg/ pkg

# build
FROM base AS build

ARG TARGETOS
ARG TARGETARCH
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 go build -o gwrp -ldflags="-s -w" cmd/gwrp/gwrp.go

# prod
FROM alpine:3.13.6 as prod

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
COPY --chown=appuser:appgroup --from=build /go/src/github.com/shteou/gwrp/gwrp /usr/bin/gwrp
EXPOSE 8080
USER appuser
CMD ["gwrp"]
