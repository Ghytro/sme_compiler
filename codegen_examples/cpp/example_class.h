#ifndef EXAMPLE_CLASS_H
#define EXAMPLE_CLASS_H

#include <istream>
#include <ostream>
#include <vector>
#include <unordered_map>
#include <sstream>

class ParseErrorException: public std::exception {
public:
    const char* what() const throw () {
        return "Incorrect format of incoming binary data, could not parse";
    }
};

class BaseSmeStruct {
public:
    virtual void FromString(const std::string& bytes) final {
        std::stringstream ss;
        ss.str(bytes);
        FromIstream(ss);
    }

    virtual std::string ToString() const final {
        std::stringstream ss;
        WriteToOstream(ss);
        return ss.str();
    }

    virtual void FromIstream(std::istream&) = 0;

    virtual void WriteToOstream(std::ostream&) const = 0;
};

class NestedStruct: public BaseSmeStruct {
public:
    NestedStruct() {
    }

    unsigned int GetField1() const {
        return field1;
    }
    void SetField1(unsigned int value) {
        field1 = value;
    }

    long long GetField2() const {
        return field2;
    }
    void SetField2(long long value) {
        field2 = value;
    }

    double GetField3() const {
        return field3;
    }
    void SetField3(double value) {
        field3 = value;
    }

    double GetField4() const {
        return field4;
    }
    void SetField4(double value) {
        field4 = value;
    }

    void FromIstream(std::istream& is) override {
        is.read(reinterpret_cast<char*>(&field1), sizeof(field1));
        is.read(reinterpret_cast<char*>(&field2), sizeof(field2));
        is.read(reinterpret_cast<char*>(&field3), sizeof(field3));
        is.read(reinterpret_cast<char*>(&field4), sizeof(field4));
    }

    void WriteToOstream(std::ostream& os) const override {
        os.write(reinterpret_cast<const char*>(&field1), sizeof(field1));
        os.write(reinterpret_cast<const char*>(&field2), sizeof(field2));
        os.write(reinterpret_cast<const char*>(&field3), sizeof(field3));
        os.write(reinterpret_cast<const char*>(&field4), sizeof(field4));
    }

private:
    uint32_t field1;
    int64_t field2;
    double field3, field4;
};

// a simple structure with fields of types:
// uint32, int64, double, double, string, list of ints and nested struct
class ExampleClass1: public BaseSmeStruct {
public:
    ExampleClass1() {

    }

    uint32_t GetField1() const {
        return field1;
    }
    void SetField1(uint32_t value) {
        field1 = value;
    }

    int64_t GetField2() const {
        return field2;
    }
    void SetField2(int64_t value) {
        field2 = value;
    }

    double GetField3() const {
        return field3;
    }
    void SetField3(double value) {
        field3 = value;
    }

    double GetField4() const {
        return field4;
    }
    void SetField4(double value) {
        field4 = value;
    }

    std::string GetField5() const {
        return field5;
    }
    void SetField5(const std::string& value) {
        field5 = value;
    }

    std::vector<uint32_t>& GetField6() {
        return field6;
    }

    NestedStruct& GetField7() {
        return field7;
    }

    std::vector<NestedStruct>& GetField8() {
        return field8;
    }

    void FromIstream(std::istream& is) {
        uint32_t currContainerSize = 0;

        //parsing primitive types
        is.read(reinterpret_cast<char*>(&field1), sizeof(field1));
        is.read(reinterpret_cast<char*>(&field2), sizeof(field2));
        is.read(reinterpret_cast<char*>(&field3), sizeof(field3));
        is.read(reinterpret_cast<char*>(&field4), sizeof(field4));

        //parsing string
        //getting size of string
        is.read(reinterpret_cast<char*>(&currContainerSize), sizeof(currContainerSize));
        //getting string itself
        field5.resize(currContainerSize);
        is.read(field5.data(), sizeof(char) * currContainerSize);

        //parsing array
        //getting size of array
        is.read(reinterpret_cast<char*>(&currContainerSize), sizeof(currContainerSize));
        //getting array itself
        field6.resize(currContainerSize);
        for (auto &x: field6) {
            is.read(reinterpret_cast<char*>(&x), sizeof(x));
        }

        //getting nested structs
        field7.FromIstream(is);

        //getting array of structs
        is.read(reinterpret_cast<char*>(&currContainerSize), sizeof(currContainerSize));
        field8.resize(currContainerSize);
        for (auto &x: field8) {
            x.FromIstream(is);
        }

        // getting map of primitive types
        // getting size of map
        is.read(reinterpret_cast<char*>(&currContainerSize), sizeof(currContainerSize));
        // getting the map itself
        for (uint32_t i = 0; i < currContainerSize; ++i) {
            std::pair<uint32_t, uint32_t> p;
            is.read(reinterpret_cast<char*>(&p.first), sizeof(p.first));
            is.read(reinterpret_cast<char*>(&p.second), sizeof(p.second));
            field9.insert(std::move(p));
        }

        // getting map of user defined structs
        // getting size of map
        is.read(reinterpret_cast<char*>(&currContainerSize), sizeof(currContainerSize));
        // getting the map itself
        for (uint32_t i = 0; i < currContainerSize; ++i) {
            std::pair<uint32_t, NestedStruct> p;
            is.read(reinterpret_cast<char*>(&p.first), sizeof(p.first));
            p.second.FromIstream(is);
            field10.insert(std::move(p));
        }
    }

    void WriteToOstream(std::ostream& os) const {
        uint32_t currContainerSize = 0;

        os.write(reinterpret_cast<const char*>(&field1), sizeof(field1));
        os.write(reinterpret_cast<const char*>(&field2), sizeof(field2));
        os.write(reinterpret_cast<const char*>(&field3), sizeof(field3));
        os.write(reinterpret_cast<const char*>(&field4), sizeof(field4));

        currContainerSize = field5.length();
        os.write(reinterpret_cast<const char*>(&currContainerSize), sizeof(currContainerSize));
        os.write(field5.data(), sizeof(char) * currContainerSize);

        currContainerSize = field6.size();
        os.write(reinterpret_cast<const char*>(&currContainerSize), sizeof(currContainerSize));
        for (const auto &x: field6) {
            os.write(reinterpret_cast<const char*>(&x), sizeof(x));
        }

        field7.WriteToOstream(os);

        currContainerSize = field8.size();
        os.write(reinterpret_cast<const char*>(&currContainerSize), sizeof(currContainerSize));
        for (const auto &x: field8) {
            x.WriteToOstream(os);
        }
    }

private:
    uint32_t field1;
    int64_t field2;
    double field3, field4;
    std::string field5;
    std::vector<uint32_t> field6;
    NestedStruct field7;
    std::vector<NestedStruct> field8;
    std::unordered_map<uint32_t, uint32_t> field9;
    std::unordered_map<uint32_t, NestedStruct> field10;
};

template<class T>
std::ostream& operator<<(std::ostream& os, const std::vector<T>& v) {
    for (const auto &x: v) {
        os << x << " ";
    }
    return os;
}

std::ostream& operator<<(std::ostream& os, const NestedStruct& ns) {
    os << "{" << ns.GetField1() << ";" << ns.GetField2() << ";" << ns.GetField3() << ";" << ns.GetField4() << "}";
    return os;
}

#endif // EXAMPLE_CLASS_H
