cmake_minimum_required(VERSION 3.16)

project(runtime)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -O0 -g3")

set(SRC runtime.hpp runtime.cpp program.cpp)

include_directories(.)

add_executable(runtime ${SRC})