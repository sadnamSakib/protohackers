# Step 1: Use the official Golang image to build the Go application
FROM golang:1.22.3-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod file if it exists
COPY go.mod ./

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o main ./main.go


# Expose the port the application runs on
EXPOSE 8080

# Command to run the application
CMD ["./main"]
