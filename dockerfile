# Build stage
FROM golang:1.23.3-alpine AS builder

WORKDIR /app

# Install specific version of swag (matches your go.mod)
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.2

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Generate Swagger docs with compatible version
RUN swag init -g cmd/main.go --output docs --parseDependency --parseInternal

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service ./cmd/main.go

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/auth-service .
COPY --from=builder /app/docs ./docs
EXPOSE 8000
CMD ["./auth-service"]