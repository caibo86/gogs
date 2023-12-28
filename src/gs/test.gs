import "gss"

// 前面的注释

// 这是一个学生
// 这还是一个学生
struct Student {
    // ID 唯一ID
    ID   int64  = 1; // 这是一个ID
    Name string = 2; // 名字
    Age  int32  = 300; // 年龄
}

// 颜色
enum Color { // 另外个注释
    // 上面的红色
    Red = 1; // 后面红色
    Green = 2;
    Blue = 3;
}

// 有一个属性
// 错误码
@gslang.Error
// 各种注释
enum ErrCode {
    OK = 0; // 成功
    Fail = 1; // 失败
    SystemErr = 10000; // 系统错误
}

// 汽车
struct Car {
    // 颜色
    VarEnum    Color                     = 1;
    VarString  string                    = 3;
    VarByte    byte                      = 5;
    VarSbyte   sbyte                     = 6;
    VarInt16   int16                     = 7;
    VarUint16  uint16                    = 8;
    VarInt32   int32                     = 9;// 整数
    VarUint32  uint32                    = 10;
    VarInt64   int64                     = 11;
    // 64位无符号整数
    VarUint64  uint64                    = 12;
    VarFloat32 float32                   = 14;
    VarFloat64 float64                   = 15;
    VarStruct  Student                   = 21;
    VarArray []int32                     = 22;
    VarStructs []Student                 = 23;
    VarBools []bool                      = 24;
    VarStrings []string                  = 25;
    VarFloat32s []float32                = 26;
    VarFloat64s []float64                = 27;
    VarEnums  []Color                     = 28;
    VarMap     map[string]string        = 29;
    VarMap1    map[string]Student       = 30;
    VarSubject gss.Subject               = 31;
    VarTeacher gss.Teacher               = 41;
    VarTeachers []gss.Teacher            = 42;
    VarSubjects []gss.Subject            = 43;
    VarMap2 map[gss.Subject]gss.Subject = 44;
    VarMap3 map[gss.Subject]gss.Teacher = 1000; // 这是一个学科老师映射
    VarBool bool                         = 4;
}

// 最后的日志