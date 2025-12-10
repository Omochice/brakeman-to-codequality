# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Copy dependency files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a \
    -o brakeman-to-codequality \
    .

# Runtime stage using distroless
FROM gcr.io/distroless/static-debian12:nonroot

LABEL maintainer="Omochice"
LABEL description="Convert Brakeman security scan results to GitLab Code Quality format"

# Copy binary from builder
COPY --from=builder /build/brakeman-to-codequality /usr/local/bin/brakeman-to-codequality

# Use nonroot user (UID 65532)
USER nonroot:nonroot

# Set entrypoint
ENTRYPOINT ["/usr/local/bin/brakeman-to-codequality"]
