# Use the official Golang image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .

# Download and install any required third-party dependencies into the container
RUN go mod download

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["go", "run", "main.go", "db.go"]
