# Use Go v1.25.5 as the base image
FROM golang:1.25.5-alpine

# Create a new user in the docker image
RUN adduser --disabled-password --gecos '' gouser

# Create a new directory for goserve files and set the path in the container
RUN mkdir -p /home/gouser/goserve

# Set the working directory in the container
WORKDIR /home/gouser/goserve

# Copy the project files into the container
COPY . .

# Set the ownership of the goserve directory to gouser
RUN chown -R gouser:gouser /home/gouser/goserve

# Switch to the gouser user
USER gouser

# Download dependencies and build the project
RUN go mod tidy
RUN go build -o build/server cmd/main.go

# Expose the server port (replace 8080 with your actual port)
EXPOSE 8080

# Command to run the server
CMD ["./build/server"]