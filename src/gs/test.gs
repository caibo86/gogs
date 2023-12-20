@gslang.Error
enum Err(uint16) {
    Success(0), // 成功
    OOS(1)      // 服务不可用
}

struct Student {
    ID   int64;
    Name string;
    Age  int16;
}

enum Color(uint32) {
    Red(1),
    Green(2),
    Blue(3)
}

struct Car {
    ID    int64;
    Name  string;
    Price float32;
    Owner Student;
    Code  Err;
    Size[4] int32;
    Drivers[] Student;
    Attrs map[string]string;
    Color Color;
}