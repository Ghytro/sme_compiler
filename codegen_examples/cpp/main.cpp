#include <iostream>
#include <numeric>
#include "example_class.h"

int main()
{
    ExampleClass1 obj;
    NestedStruct defaultNestedStruct;
    defaultNestedStruct.SetField1(4);
    defaultNestedStruct.SetField2(6);
    defaultNestedStruct.SetField3(1.5);
    defaultNestedStruct.SetField4(4.8);

    obj.SetField1(4);
    obj.SetField2(6);
    obj.SetField3(1.5);
    obj.SetField4(4.8);
    obj.SetField5("abacaba");
    obj.GetField6().resize(5);
    std::iota(obj.GetField6().begin(), obj.GetField6().end(), 0);
    obj.GetField7() = defaultNestedStruct;
    obj.GetField8().push_back(defaultNestedStruct);

    std::cout << obj.GetField1() << "\n"
              << obj.GetField2() << "\n"
              << obj.GetField3() << "\n"
              << obj.GetField4() << "\n"
              << obj.GetField5() << "\n"
              << obj.GetField6() << "\n"
              << obj.GetField7() << "\n"
              << obj.GetField8() << "\n";

    std::string marshaled = obj.ToString();

    ExampleClass1 unmarshaled;
    unmarshaled.FromString(marshaled);
    std::cout << unmarshaled.GetField1() << "\n"
              << unmarshaled.GetField2() << "\n"
              << unmarshaled.GetField3() << "\n"
              << unmarshaled.GetField4() << "\n"
              << unmarshaled.GetField5() << "\n"
              << unmarshaled.GetField6() << "\n"
              << unmarshaled.GetField7() << "\n"
              << unmarshaled.GetField8();
}
