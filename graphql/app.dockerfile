# -----------------------
# 1️⃣ Build Stage
# -----------------------
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc g++ make ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy module manifests and download deps
COPY go.mod go.sum ./
RUN go mod download

# Copy source code needed for graphql build
COPY account account
COPY catalog catalog
COPY order order
COPY graphql graphql

# Build the graphql gateway binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./graphql

# -----------------------
# 2️⃣ Runtime Stage
# -----------------------
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /usr/local/bin

# Copy compiled binary from builder stage
COPY --from=builder /app/app .

EXPOSE 8080

ENTRYPOINT ["./app"]
