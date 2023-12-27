import "gss"

struct Student {
    ID   int64 = 1;
    Name string = 2;
    Age  int32 = 3;
}

enum Color {
    Red(1);
    Green(2);
    Blue(3);
}

struct Car {
    VarEnum Color = 2;
    VarString string = 3;
    VarBool bool = 4;
    VarByte byte = 5;
    VarSbyte sbyte = 6;
    VarInt16 int16 = 7;
    VarUint16 uint16 = 8;
    VarInt32 int32 = 9;
    VarUint32 uint32 = 10;
    VarInt64 int64 = 11;
    VarUint64 uint64 = 12;
    VarFloat32 float32 = 14;
    VarFloat64 float64 = 15;
    VarStruct Student = 21;
    VarArray []int32 = 22;
    VarStructs []Student = 23;
    VarBools []bool =24;
    VarStrings []string = 25;
    VarFloat32s []float32 =26;
    VarFloat64s []float64 = 27;
    VarEnums []Color = 28;
    VarMap map[string]string = 29;
    VarMap1 map[string]Student= 30;
    VarSubject gss.Subject =31;
    VarTeacher gss.Teacher = 41;
    VarTeachers []gss.Teacher = 42;
    VarSubjects []gss.Subject = 43;
    VarMap2 map[gss.Subject]gss.Subject = 44;
    VarMap3 map[gss.Subject]gss.Teacher = 45;
}

