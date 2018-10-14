FROM golang:1.11-alpine AS build-base

RUN apk add bash ca-certificates git gcc g++ libc-dev
WORKDIR /go/src/github.com/liam-j-bennett/prepaidcard
ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .
RUN go mod download

FROM build-base AS build-env

COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app

FROM alpine:latest

COPY --from=build-env /go/src/github.com/liam-j-bennett/prepaidcard/app .
CMD ["/app"]