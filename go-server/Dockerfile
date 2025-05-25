# 1. Use the official Go image as the builder
FROM golang:1.24.3 AS builder

# 2. Set the working directory inside the container
WORKDIR /app

# 3. Copy go.mod and go.sum first, then download dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# 4. Copy the rest of the application
COPY . .

# 5. Build the Go app (assumes main.go or cmd/main.go exists)
RUN go build -o server .

# 6. Use a minimal base image for final container (security + size)
FROM alpine:latest

# 7. Set the working directory again
WORKDIR /app

# 8. Copy the compiled Go binary from builder
COPY --from=builder /app/server .

# 9. Expose the app port (optional)
EXPOSE 3000

# 10. Run the binary
CMD ["./server"]

