FROM golang:1.19 AS development
WORKDIR /go/src/github.com/escalopa/fingo
COPY ./contact ./contact
COPY ./pkg ./pkg
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/cespare/reflex@latest
CMD reflex -sr '\.go$' go run ./contact/cmd/.

FROM golang:alpine AS builder
WORKDIR /go/src/github.com/escalopa/fingo
COPY ./contact ./contact
COPY ./pkg ./pkg
COPY go.mod go.sum ./
RUN go build -o /go/bin/contact ./contact/cmd/.

FROM alpine:latest AS production
RUN apk add --no-cache tzdata
COPY --from=builder /go/bin/contact /go/bin/fingo-contact
COPY ./contact/app.env /contact/app.env
COPY ./certs /certs
ENTRYPOINT ["/go/bin/fingo-contact"]
