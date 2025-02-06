# Use the official Golang image to build the app
FROM golang:1.20 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules file and download dependencies
COPY go.mod ./
RUN go mod tidy

# Copy the rest of the application
COPY . .

# Build the Go application
RUN go build -o receipt-processor

# Use a minimal base image to run the application
FROM alpine:latest
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/receipt-processor .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./receipt-processor"]
