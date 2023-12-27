// -------------------------------------------
// @file      : enum.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/16 下午5:51
// -------------------------------------------

package ast

// EnumVal 枚举值 指一个枚举括号中的单个枚举值
type EnumVal struct {
	BaseExpr       // 内嵌基本表达式实现
	Value    int32 // 枚举值节点对应的实际枚举数值
}

// Enum 枚举表达式
type Enum struct {
	BaseExpr                     // 内嵌基本表达式实现
	Values   map[string]*EnumVal // 枚举值字典
	Default  *EnumVal            // 入口枚举值
}

// NewEnum 在代码节点内新建枚举节点 所属代码节点为此代码节点
func (script *Script) NewEnum(name string) *Enum {
	enum := &Enum{
		Values: make(map[string]*EnumVal),
	}
	// 初始化枚举表达式
	enum.Init(name, script)
	// 设置父节点为此代码节点
	enum.SetParent(script)
	return enum
}

// NewEnumVal 在枚举内新建一个枚举值
func (enum *Enum) NewEnumVal(name string, val int32) (*EnumVal, bool) {
	// 检查枚举表达式内是否已有同名枚举值 有则直接返回
	enumVal, ok := enum.Values[name]
	if ok {
		return enumVal, false
	}
	// 新建枚举值
	enumVal = &EnumVal{
		Value: val,
	}
	// 初始化枚举值,所属代码节点为枚举表达式所属代码节点
	enumVal.Init(name, enum.Script())
	// 将枚举值加入到枚举表达式的枚举值字典
	enum.Values[name] = enumVal
	// 如果枚举表达式的入口枚举值为空 则将此枚举值设置为入口枚举值
	if enum.Default == nil {
		enum.Default = enumVal
	}
	return enumVal, true
}
