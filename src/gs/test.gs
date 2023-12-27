import "gss"

struct Student {
    ID   int64;
    Name string;
    Age  int32;
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
    VarSbyte sbyte;
    VarInt16 int16;
    VarUint16 uint16;
    VarInt32 int32;
    VarUint32 uint32;
    VarInt64 int64;
    VarUint64 uint64;
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
    VarMap map[string]string;
    VarMap1 map[string]Student;
    VarSubject gss.Subject;
    VarTeacher gss.Teacher;
    VarTeachers []gss.Teacher;
    VarSubjects []gss.Subject;
    VarMap2 map[gss.Subject]gss.Subject;
    VarMap3 map[gss.Subject]gss.Teacher;
}

