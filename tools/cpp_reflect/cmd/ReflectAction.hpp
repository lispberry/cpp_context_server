#pragma once

#pragma once

#include <memory>
#include <utility>

#include "clang/AST/RecursiveASTVisitor.h"
#include "clang/Basic/SourceManager.h"
#include "clang/Frontend/CompilerInstance.h"
#include "clang/Frontend/FrontendAction.h"
#include "clang/Tooling/Tooling.h"

#include "cmd/ReflectAstConsumer.hpp"

class ReflectAction : public clang::ASTFrontendAction
{
public:
  std::unique_ptr<clang::ASTConsumer> CreateASTConsumer(clang::CompilerInstance & compilerInstance,
                                                        clang::StringRef) final
  {
    compilerInstance.getDiagnostics().setClient(new clang::IgnoringDiagConsumer());
    return std::make_unique<ReflectAstConsumer>(compilerInstance.getSourceManager());
  }

private:
};