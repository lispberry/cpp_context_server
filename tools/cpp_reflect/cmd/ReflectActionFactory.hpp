#pragma once

#include "clang/Frontend/FrontendAction.h"
#include "clang/Tooling/Tooling.h"

#include "cmd/ReflectAction.hpp"

#include <memory>

class ReflectActionFactory : public clang::tooling::FrontendActionFactory
{
public:
  std::unique_ptr<clang::FrontendAction> create() final
  {
    return std::make_unique<ReflectAction>();
  }

private:
};

inline std::unique_ptr<ReflectActionFactory> newReflectActionFactory()
{
  return std::make_unique<ReflectActionFactory>();
}