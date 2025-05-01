# ======== Build Stage ========
FROM golang:1.22 AS builder

WORKDIR /app

# Copy go.mod and download dependencies
COPY go.mod ./
RUN go mod download

# Copy source files
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp main.go

# ======== Final Stage ========
FROM alpine:latest

WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/myapp .

# Expose app port
EXPOSE 8080

# Run the binary
CMD ["./myapp"]
