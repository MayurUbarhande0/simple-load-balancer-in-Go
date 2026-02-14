# Build stage
FROM golang:1.25.3-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o lb .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/lb .
COPY --from=builder /app/config.example.json ./config.json

# Expose the load balancer port
EXPOSE 8080

# Run the load balancer
CMD ["./lb"]
