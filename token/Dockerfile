FROM golang:1.19 AS development
WORKDIR /go/src/github.com/escalopa/fingo
COPY ./token ./token
COPY ./certs ./certs
COPY ./pkg ./pkg
COPY ./pb ./pb
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/cespare/reflex@latest
CMD reflex -sr '\.go$' go run ./token/cmd/.

FROM golang:alpine AS builder
WORKDIR /go/src/github.com/escalopa/fingo
COPY ./token ./token
COPY ./pkg ./pkg
COPY ./pb ./pb
COPY go.mod go.sum ./
RUN go build -o /go/bin/token ./token/cmd/.

FROM alpine:latest AS production
RUN apk add --no-cache tzdata
COPY --from=builder /go/bin/token /go/bin/fingo-token
COPY ./token/app.env /token/app.env
COPY ./certs /certs
ENTRYPOINT ["/go/bin/fingo-token"]
