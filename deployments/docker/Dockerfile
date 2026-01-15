# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o notification-service ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/notification-service .

# Copy config files and Firebase credentials
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/corechain-e1321-firebase-adminsdk-fbsvc-fc8bac45e8.json .

# Add metadata labels
LABEL maintainer="CoreChain Team"
LABEL version="1.0.0"
LABEL description="CoreChain Notification Service - Kafka to FCM notification delivery"

# Expose health check port (optional)
EXPOSE 8080

CMD ["./notification-service"]
