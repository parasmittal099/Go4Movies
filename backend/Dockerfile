# Use the official Go image
FROM golang:1.26.0-alpine

# Install C compiler for CGO (required by gorm.io/driver/sqlite)
RUN apk add --no-cache gcc musl-dev

# Set the working directory
WORKDIR /app

# Copy go mod files first to leverage Docker caching
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Create data directory for SQLite
RUN mkdir -p /app/data

# Expose the port the app runs on
EXPOSE 8080

# Run the application
CMD ["sh", "-c", "rm -f /app/data/*.db /app/data/*.db-shm /app/data/*.db-wal && go run main.go"]
