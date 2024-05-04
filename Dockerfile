# Build stage
FROM golang:1.22-alpine AS build

# Set the working directory
WORKDIR /app

# Download necessary Go modules
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the source code
COPY cmd/ ./cmd/
COPY pkg/ ./pkg/

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server/main.go

# Run stage
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the binary from the build stage
COPY --from=build /app/main .

# Expose port 8080
EXPOSE 8080

# Set the command to run
CMD ["./main"]