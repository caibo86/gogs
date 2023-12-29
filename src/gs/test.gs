import "gss"

// 前面的注释
// 最后的日志

// 颜色
enum Color {
	Red   = 1; // 另外个注释 上面的红色 后面红色
	Green = 2;
	Blue  = 3;
}

// 有一个属性
// 错误码
// 各种注释
@gslang.Error
enum ErrCode {
	OK        = 0;     // 成功
	Fail      = 1;     // 失败
	SystemErr = 10000; // 系统错误
}

struct Phone {
	Number      string = 1; // 前面的注释 这是一个手机号
	CountryCode int32  = 2; // 这是一个国家代码
}

// 这是一个学生
// 这还是一个学生
// +k8s:deepcopy-gen=true
struct Student {
	ID    int64  = 1;   // ID 唯一ID 这是一个ID
	Name  string = 2;   // 名字
	Phone Phone  = 4;   // 手机
	Age   int32  = 300; // 年龄
}

struct Table {
}

@gslang.AttrUsage(Target:gslang.AttrTarget.Struct)
table Logo {
	ID    int32   = 1;
	Name  string  = 2;
	Color Color   = 3;
	B     bool    = 4;
	F     float32 = 5;
	D     float64 = 6;
	Num1  byte    = 7;
	Num2  int8    = 8;
	Num3  int16   = 9;
	Num4  uint16  = 10;
	Num5  int32   = 11;
	Num6  uint32  = 12;
	Num7  int64   = 13;
	Num8  uint64  = 14;
	Num9  uint8   = 15;
}

// 汽车
// +k8s:deepcopy-gen=true
@Logo(ID:100, Name:"logo", Color:Color.Blue, B:true, F:3.1e-3, D:-9.99, Num1:1, Num2:2, Num3:3, Num4:4, Num5:0xFF, Num6:0x12AB, Num7:7, Num8:8, Num9:9)
struct Car {
	VarEnum         Color                       = 1;    // 颜色
	VarUint8        uint8                       = 2;
	VarString       string                      = 3;
	VarBool         bool                        = 4;
	VarByte         byte                        = 5;
	VarInt8         int8                        = 6;
	VarInt16        int16                       = 7;
	VarUint16       uint16                      = 8;
	VarInt32        int32                       = 9;    // 整数
	VarUint32       uint32                      = 10;
	VarInt64        int64                       = 11;
	VarUint64       uint64                      = 12;   // 64位无符号整数
	VarFloat32      float32                     = 14;
	VarFloat64      float64                     = 15;
	VarStruct       Student                     = 21;
	VarList         []int32                     = 22;
	VarStructs      []Student                   = 23;
	VarBools        []bool                      = 24;
	VarStrings      []string                    = 25;
	VarFloat32s     []float32                   = 26;
	VarFloat64s     []float64                   = 27;
	VarEnums        []Color                     = 28;
	VarMap          map[string]string           = 29;
	VarMap1         map[string]Student          = 30;
	VarSubject      gss.Subject                 = 31;
	VarTeacher      gss.Teacher                 = 41;
	VarTeachers     []gss.Teacher               = 42;
	VarSubjects     []gss.Subject               = 43;
	VarMap2         map[gss.Subject]gss.Subject = 44;
	VarArray        [3]int32                    = 45;
	VarStructArray  [3]gss.Teacher              = 46;
	VarEnumArray    [4]gss.Subject              = 47;
	VarStructArray1 [4]Student                  = 48;
	VarStructSlice  []Table                     = 49;
	VarStructMap    map[string]Table            = 50;
	VarData         []byte                      = 51;
	VarBytes        bytes                       = 52;
	VarArrayBytes   [10]byte                    = 53;
	VarMap3         map[gss.Subject]gss.Teacher = 1000; // 这是一个学科老师映射
}

// 游戏服
service GameServer {
     GetServerTime() -> (int64);
}