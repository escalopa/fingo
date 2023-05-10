FROM golang:1.19 AS development
WORKDIR /go/src/github.com/escalopa/fingo
COPY ./auth ./auth
COPY ./certs ./certs
COPY ./pkg ./pkg
COPY ./pb ./pb
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/cespare/reflex@latest
CMD reflex -sr '\.go$' go run ./auth/cmd/.

FROM golang:alpine AS builder
WORKDIR /go/src/github.com/escalopa/fingo
COPY ./auth ./auth
COPY ./pkg ./pkg
COPY ./pb ./pb
COPY go.mod go.sum ./
RUN go build -o /go/bin/auth ./auth/cmd/.

FROM alpine:latest AS production
RUN apk add --no-cache tzdata
COPY --from=builder /go/bin/auth /go/bin/fingo-auth
COPY ./auth/internal/adapters/db/postgres/migrations /migrations
COPY ./auth/app.env /auth/app.env
COPY ./certs /certs
ENTRYPOINT ["/go/bin/fingo-auth"]
