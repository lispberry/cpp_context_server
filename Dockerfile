FROM ubuntu:latest as cpp_reflect
RUN apt-get update && apt-get install -y \
    build-essential \
    gdb \
    cmake \
    libclang-dev \
    clang \
    llvm \
    golang-go \
    ninja-build && \
    rm -rf /var/lib/apt/lists/* \

ENV CGO_ENABLED=1

WORKDIR /app

COPY . .

RUN CC=clang CXX=clang++ cmake -S ./tools/cpp_reflect/ -B ./tools/cpp_reflect/release -DCMAKE_BUILD_TYPE=Release -GNinja .
RUN cd ./tools/cpp_reflect/release && ninja -j4 && ninja install

RUN go mod download

CMD ["go", "run", "cmd/main.go"]