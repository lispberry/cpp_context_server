#pragma once

#include <string>
#include <variant>
#include <sstream>
#include <iostream>
#include "json.hpp"
using json = nlohmann::json;

struct List {
  int val{0};
  List *n{nullptr};
};

using Execution = std::string;

inline void Write(json const & array)
{
  std::cout << array.dump() << std::endl;
}

namespace internal {
inline json createRawOp(const std::string &kind, json data) {
  json rawOp;
  rawOp["kind"] = kind;
  rawOp["data"] = data;
  return rawOp;
}

// Function to create a JSON object of type NewListPointer
inline json createNewListPointer(const std::string &name) {
  json newListPointer;
  newListPointer["name"] = name;
  return createRawOp("NewListPointer", newListPointer);
}

// Function to create a JSON object of type SetListPointerValue
inline json createSetListPointerValue(const std::string &name, const std::string &address) {
  json setListPointerValue;
  setListPointerValue["name"] = name;
  setListPointerValue["address"] = address;
  return createRawOp("SetListPointerValue", setListPointerValue);
}

// Function to create a JSON object of type NewListNode
inline json createNewListNode(const std::string &address) {
  json newListNode;
  newListNode["address"] = address;
  return createRawOp("NewListNode", newListNode);
}

// Function to create a JSON object of type SetListNodeNext
inline json createSetListNodeNext(const std::string &address, const std::string &next) {
  json setListNodeNext;
  setListNodeNext["address"] = address;
  setListNodeNext["next"] = next;
  return createRawOp("SetListNodeNext", setListNodeNext);
}

// Function to create a JSON object of type SetListNodeValue
inline json createSetListNodeValue(const std::string &address, const std::string &value) {
  json setListNodeValue;
  setListNodeValue["address"] = address;
  setListNodeValue["value"] = value;
  return createRawOp("SetListNodeValue", setListNodeValue);
}
}

template <typename T>
class Pointer
{
public:
  explicit Pointer(std::string const &name)
    : m_name(name)
    , m_data(nullptr)
  {
    Write(json::array({
      internal::createNewListPointer(name)
    }));
  }

  explicit Pointer(T* data)
    : m_name("")
    , m_data(data)
  {
      Write(json::array({
        internal::createNewListNode(Address())
      }));
  }

  Pointer(Pointer const & that) {
    m_data = that.m_data;
  }
  Pointer(Pointer && that)  noexcept {
    m_data = that.m_data;
  }
  Pointer & operator=(const Pointer &that) {
    m_data = that.m_data;
  }
  Pointer & operator=(Pointer &&that) {
    m_data = that.m_data;
  }

  Pointer & operator=(T * data) {
    m_data = data;
    return *this;
  }

  operator T*() const {
    return m_data;
  }

  T& operator*() const noexcept {
    return *m_data;
  }

  T* operator->() const noexcept {
    return m_data;
  }

  void SetName(const std::string & name) {
    m_name = name;
  }

  void SetAddress(T * addr) {
    m_data = addr;
  }

  std::string Address() const {
    std::stringstream ss;
    ss << m_data;
    return ss.str();
  }

  std::string Name() const {
    return m_name;
  }
private:
  std::string m_name;
  T * m_data;
};

template class Pointer<List>;

template <typename T>
inline std::string FromAddress(T * ptr)
{
  std::stringstream ss;
  ss << ptr;
  return ss.str();
}

inline void SetListPointerValue(Pointer<List> & ptr) {
  ptr.SetAddress(new List);

  Write(json::array({
    internal::createNewListNode(ptr.Address()),
    internal::createSetListPointerValue(ptr.Name(), ptr.Address())
  }));
}

inline void SetListPointerValue(Pointer<List> & ptr, List *addr)
{
  ptr.SetAddress(addr);
  Write(json::array({
                      internal::createSetListPointerValue(ptr.Name(), ptr.Address())
                    }));
}

inline void SetListPointerValue(Pointer<List> & ptr, Pointer<List> & that)
{
  ptr.SetAddress(that);
  Write(json::array({
    internal::createSetListPointerValue(ptr.Name(), ptr.Address())
  }));
}

inline void SetListNodeNext(Pointer<List> &ptr, Pointer<List> & that)
{
  ptr->n = that;
  Write(json::array({
    internal::createSetListNodeNext(ptr.Address(), that.Address())
  }));
}

inline void SetListNodeNext(Pointer<List> & ptr, std::nullptr_t null)
{
  ptr->n = nullptr;
  Write(json::array({
    internal::createSetListNodeNext(ptr.Address(), "0x0")
  }));
}

inline void SetListNodeValue(Pointer<List> & ptr, int value)
{
  ptr->val = value;
  Write(json::array({
    internal::createSetListNodeValue(ptr.Address(), std::to_string(value))
  }));
}

inline Pointer<List>& unmove(Pointer<List>&& t) { return t; }
inline Pointer<List>& unmove(Pointer<List>& t) { return t; }

extern Pointer<List> DefaultList;

inline Pointer<List> CreatePointer(const char * cname) {
    std::string name = cname;

    auto cp = DefaultList;
    cp.SetName(name);
    Write(json::array({
      internal::createNewListPointer(name),
      internal::createSetListPointerValue(name, cp.Address())
    }));

    return cp;
}