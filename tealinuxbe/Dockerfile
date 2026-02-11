# ============================================
# Stage 1: Build
# ============================================
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git (needed for some go modules)
RUN apk add --no-cache git

# Copy go.mod and go.sum first for better layer caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o tealinux-api ./cmd/main.go

# ============================================
# Stage 2: Production
# ============================================
FROM alpine:3.19 AS runtime

WORKDIR /app

# Install ca-certificates for HTTPS requests
RUN apk add --no-cache ca-certificates tzdata

# Copy the binary from builder
COPY --from=builder /app/tealinux-api .

# Set environment
ENV GIN_MODE=release

EXPOSE 3000

# Run the binary
CMD ["./tealinux-api"]
