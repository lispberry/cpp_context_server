project(cpp_reflect_cmd)

set(SRC
  main.cpp
  ReflectActionFactory.hpp
  ReflectAction.hpp
  ReflectAstConsumer.hpp)

add_executable(${PROJECT_NAME} ${SRC})

target_include_directories(${PROJECT_NAME} PUBLIC ${LLVM_INCLUDE_DIRS} .)
llvm_map_components_to_libnames(${PROJECT_NAME} support core irreader)
target_link_libraries(${PROJECT_NAME}
  ${llvm_libs}
  nlohmann_json::nlohmann_json
  clangAST
  clangBasic
  clangDriver
  clangFrontend
  clangRewriteFrontend
  clangSerialization
  clangStaticAnalyzerFrontend
  clangTooling
  clangToolingSyntax)

install(TARGETS ${PROJECT_NAME} DESTINATION /usr/local/bin/)