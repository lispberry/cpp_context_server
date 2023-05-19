#include "runtime.hpp"

Pointer<List> DefaultList(new List);

int main() {
  unmove(DefaultList);
  CreatePointer("Default");

  return 0;
}