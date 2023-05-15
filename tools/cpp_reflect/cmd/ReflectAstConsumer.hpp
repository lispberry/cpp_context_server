#pragma once

#include <iostream>

#include "clang/AST/ASTConsumer.h"
#include "clang/AST/RecursiveASTVisitor.h"
#include "clang/Frontend/CompilerInstance.h"
#include "clang/Basic/SourceManager.h"

#include "cmd/ast.hpp"

namespace
{
class AstVisitor : public clang::RecursiveASTVisitor<AstVisitor>
{
public:
  explicit AstVisitor(ast::File & file) : m_file(file) {}

  bool VisitFunctionDecl(clang::FunctionDecl *decl)
  {
    std::vector<ast::Argument> arguments;
    for (int i = 0; i < decl->getNumParams(); i++) {
      auto const & param = decl->getParamDecl(i);
      arguments.emplace_back(ast::Argument{
        .name = param->getNameAsString(),
        .type = param->getType().getAsString()
      });
    }

    auto loc = decl->getLocation();
    auto & sm = decl->getASTContext().getSourceManager();
    int lineNumber = sm.getPresumedLineNumber(loc);

    auto & fun = m_file.functions.emplace_back(ast::Function{
      .name = decl->getNameAsString(),
      .arguments = std::move(arguments),
      .line_number = lineNumber
    });

    // Process function declarations
    std::string functionName = decl->getNameAsString();

    return true;
  }

private:
  ast::File & m_file;
};
}

class ReflectAstConsumer : public clang::ASTConsumer
{
public:
  explicit ReflectAstConsumer(clang::SourceManager & sourceManager)
    : m_sourceManager(sourceManager) {}

  void HandleTranslationUnit(clang::ASTContext & context) final
  {
    AstVisitor visitor(file);
    visitor.TraverseDecl(context.getTranslationUnitDecl());

    std::cout << ast::SerializeToJson(file);
  }

private:
  ast::File file;
  clang::SourceManager & m_sourceManager;
};