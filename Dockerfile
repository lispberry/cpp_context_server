FROM golang:latest

# Install GDB
RUN apt-get update && \
    apt-get install -y gdb gcc && \
    rm -rf /var/lib/apt/lists/*

# Enable CGO
ENV CGO_ENABLED=1

# Set the working directory
WORKDIR /app

# Copy the source code to the container
COPY . .

RUN go mod download

RUN go test github.com/lispberry/viz-service/pkg/evaluation

# Start the application
CMD ["go", "run", "cmd/main.go"]