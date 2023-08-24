# Use an official Go runtime as the base image
FROM golang:1.19

# Set the working directory inside the container
WORKDIR /app

# Copy the Go source code into the container's workspace
# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./ 
RUN go mod download && go mod verify

COPY ./app/*.go ./

# Build the Go application
RUN go build -o app

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./app"]