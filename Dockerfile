FROM golang:1.25-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o brakeman-to-codequality .

FROM alpine:3

COPY --from=builder /build/brakeman-to-codequality /usr/local/bin/brakeman-to-codequality
