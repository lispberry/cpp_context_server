#include "clang/Frontend/FrontendActions.h"
#include "clang/Tooling/CommonOptionsParser.h"
#include "clang/Tooling/Tooling.h"
#include "llvm/Support/CommandLine.h"

#include "cmd/ReflectActionFactory.hpp"

using namespace clang::tooling;
using namespace llvm;

static llvm::cl::OptionCategory ReflectCategory("cpp_reflect options");

static cl::extrahelp CommonHelp(CommonOptionsParser::HelpMessage);

int main(int argc, char const ** argv)
{
  auto ExpectedParser = CommonOptionsParser::create(argc, argv, ReflectCategory);
  if (!ExpectedParser)
  {
    // Fail gracefully for unsupported options.
    llvm::errs() << ExpectedParser.takeError();
    return 1;
  }
  CommonOptionsParser & OptionsParser = ExpectedParser.get();
  ClangTool Tool(OptionsParser.getCompilations(), OptionsParser.getSourcePathList());

  return Tool.run(newReflectActionFactory().get());
}