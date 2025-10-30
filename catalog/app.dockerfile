# -----------------------
# 1️⃣ Build Stage
# -----------------------
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc g++ make ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
# COPY . .
COPY catalog catalog

# (Optional) If you have vendor directory:
# RUN go mod vendor

# Build the binary (disable CGO for static binary)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./catalog/cmd/catalog

# -----------------------
# 2️⃣ Runtime Stage
# -----------------------
FROM alpine:3.20

# Set timezone and install CA certificates
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /usr/local/bin

# Copy binary from builder
COPY --from=builder /app/app .

# Non-root user (security best practice)
# RUN adduser -D appuser
# USER appuser

# Expose the app port
EXPOSE 8080

# Start the application
ENTRYPOINT ["./app"]
