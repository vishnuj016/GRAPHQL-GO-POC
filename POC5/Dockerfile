# Use an official Go runtime as the base image
FROM golang:1.20

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to cache dependencies
COPY go.mod go.sum ./

# Download and install Go dependencies
RUN go mod download

# Copy your application code into the container
COPY . .

# Build your Go application
RUN go build -o gqlgopoc ./server.go  # Modify the path to your main Go file

# Expose the port your application listens on (if applicable)
EXPOSE 8080

# Define the command to run when the container starts
CMD ["./gqlgopoc"]
