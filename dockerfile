# Use an official Go image as the base
FROM golang:1.22.2

# Set the working directory inside the container
WORKDIR /app

# Copy only go.mod and go.sum first for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Install loclx
RUN apt-get update && apt-get install -y snapd && snap install loclx

# Copy the entire project into the container
COPY . .

# Build the main Go application
RUN go build -o main main.go

# Command to run the application
CMD ["./main"]
