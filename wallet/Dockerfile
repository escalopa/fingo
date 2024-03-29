FROM golang:1.19 AS development
WORKDIR /go/src/github.com/escalopa/fingo
COPY ./wallet ./wallet
COPY ./certs ./certs
COPY ./pkg ./pkg
COPY ./pb ./pb
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/cespare/reflex@latest
CMD reflex -sr '\.go$' go run ./wallet/cmd/.

FROM golang:alpine AS builder
WORKDIR /go/src/github.com/escalopa/fingo
COPY ./wallet ./wallet
COPY ./pkg ./pkg
COPY ./pb ./pb
COPY go.mod go.sum ./
RUN go build -o /go/bin/wallet ./wallet/cmd/.

FROM alpine:latest AS production
RUN apk add --no-cache tzdata
COPY --from=builder /go/bin/wallet /go/bin/fingo-wallet
COPY ./wallet/internal/adapters/db/sql/migrations /migrations
COPY ./wallet/app.env /wallet/app.env
COPY ./certs /certs
ENTRYPOINT ["/go/bin/fingo-wallet"]
