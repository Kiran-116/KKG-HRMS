FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
ENV GOTOOLCHAIN=auto
RUN go mod download

# Copy source code
COPY . .

# Build the application
ENV CGO_ENABLED=0
RUN go build -ldflags="-s -w" -o hrms .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/hrms .
COPY --from=builder /app/migrations ./migrations

# Create uploads directory
RUN mkdir -p uploads

EXPOSE 8080

CMD ["./hrms"]
