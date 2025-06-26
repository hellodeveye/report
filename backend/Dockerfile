# Stage 1: Build the application
FROM golang:1.19-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Set go proxy
RUN go env -w GOPROXY=https://goproxy.cn,direct

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
# -o report-backend: specifies the output file name
# -ldflags="-w -s": strips debugging information, reducing the binary size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o report-backend .

# Stage 2: Create the final, lightweight image
FROM alpine:latest

WORKDIR /root/

# Copy the pre-built binary from the previous stage
COPY --from=builder /app/report-backend .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./report-backend"] 