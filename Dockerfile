FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /brakeman-to-codequality .

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /brakeman-to-codequality /usr/local/bin/brakeman-to-codequality

ENTRYPOINT ["brakeman-to-codequality"]
