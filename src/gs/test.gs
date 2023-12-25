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
    VarEnum Color;
    VarString string;
    VarBool bool;
    VarByte byte;
    VarSByte sbyte;
    VarInt16 int16;
    VarUInt16 uint16;
    VarInt32 int32;
    VarUInt32 uint32;
    VarInt64 int64;
    VarUInt64 uint64;
    VarFloat32 float32;
    VarFloat64 float64;
    VarStruct Student;
    VarArray []int32;
    VarStructs []Student;
    VarBools []bool;
    VarStrings []string;
    VarFloat32s []float32;
    VarFloat64s []float64;
    VarEnums []Color;
}

