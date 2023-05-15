FROM ubuntu:latest as cpp_reflect

# Install necessary dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    cmake \
    libclang-dev \
    clang \
    llvm \
    ninja-build

# Set the working directory to /tools
WORKDIR /tools

# Copy the CMakeLists.txt and source code to the container
COPY ./tools/cpp_reflect/ .

# Build the program using CMake and Ninja
RUN CC=clang CXX=clang++ cmake -DCMAKE_BUILD_TYPE=Release -GNinja .
RUN ninja -j4

FROM golang:latest

# Install GDB
RUN apt-get update && \
    apt-get install -y gdb gcc && \
    rm -rf /var/lib/apt/lists/*

# Enable CGO
ENV CGO_ENABLED=1

COPY --from=cpp_reflect /tools/cmd/cpp_reflect_cmd /tools/cpp_reflect

# Set the working directory
WORKDIR /app

# Copy the source code to the container
COPY . .

RUN go mod download

RUN go test github.com/lispberry/viz-service/pkg/semantic
RUN go test github.com/lispberry/viz-service/pkg/evaluation

# Start the application
CMD ["go", "run", "cmd/main.go"]