cmake_minimum_required(VERSION 3.12)

project(cpp_reflect)

set(CMAKE_CXX_STANDARD 20)

if (${CMAKE_SYSTEM_NAME} MATCHES "Darwin")
  set(MACOSX TRUE)
endif()

if (${MACOSX})
  list(APPEND CMAKE_INCLUDE_PATH /opt/homebrew/include)
  list(APPEND CMAKE_LIBRARY_PATH /opt/homebrew/lib)
endif()

find_package(LLVM  REQUIRED CONFIG)
find_package(Clang REQUIRED CONFIG)

include_directories(.)

add_subdirectory(3party/json)

add_subdirectory(cmd)