FROM golang:1.24.4

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./
RUN go mod tidy

# Copy all other files
COPY . .

# Build the Go binary
RUN go build -o app .

# Run the binary
CMD ["./app"]
