// 属性的目标枚举
enum AttrTarget {
    Package = 1;
    Script = 2;
    Table = 4;
    Struct = 8;
    Enum = 16;
    EnumVal = 32;
    Field = 64;
    Contract = 128;
    Method = 256;
    Return = 512;
    Param = 1024;
}

// 用于标注属性目标类型的内置Table
@AttrUsage(Target:AttrTarget.Table)
table AttrUsage {
    // 用于反转查找目标类型
    Target AttrTarget = 1;
}

// 内置Struct类型用于标注Table是一个Struct
@AttrUsage(Target:AttrTarget.Struct)
table Struct {}


// 内置类型标注Enum是一个错误类型声明
@AttrUsage(Target:AttrTarget.Enum)
table Error {}
