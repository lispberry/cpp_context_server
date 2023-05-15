#pragma once

#include <vector>
#include <string>
#include <nlohmann/json.hpp>

namespace ast
{
using json = nlohmann::json;

struct Argument {
  std::string name;
  std::string type;
};

struct Function {
  std::string name;
  std::vector<Argument> arguments;
  int line_number;
};

struct File {
  std::vector<Function> functions;
};

std::string SerializeToJson(const File& file) {
  json j;
  j["functions"] = json::array();

  for (const auto& function : file.functions) {
    json jFunction;
    jFunction["name"] = function.name;
    jFunction["arguments"] = json::array();

    for (const auto& argument : function.arguments) {
      json jArgument;
      jArgument["name"] = argument.name;
      jArgument["type"] = argument.type;
      jFunction["arguments"].push_back(jArgument);
    }

    jFunction["line_number"] = function.line_number;
    j["functions"].push_back(jFunction);
  }

  return j.dump();
}
}