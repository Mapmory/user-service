# Use the official Golang image
FROM golang:1.23

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Set the working directory to /app/api
# WORKDIR /app/api

# Build the application
RUN go build -o main .

# Expose port 8000
EXPOSE 8000

# Run the executable
CMD ["./main"]