FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0

ENV GOOS linux

RUN apk update --no-cache

WORKDIR /build

ADD go.mod .

ADD go.sum .

RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o main ./cmd/main.go

FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /build/main ./

COPY .env .

RUN mkdir "config"

COPY config/config.yaml ./config

RUN mkdir -p "db/migration"

COPY db/migrations ./db/migrations

CMD ["./main"]